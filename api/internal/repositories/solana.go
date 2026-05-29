package repositories

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/odingaval/veryfy/api/internal/config"
	"github.com/odingaval/veryfy/api/internal/models"
)

var (
	ErrSolanaIntegrationPending = errors.New("solana integration pending")
	ErrInvalidSolanaConfig      = errors.New("invalid solana config")
)

const (
	anchorDiscriminatorSize = 8
	licenseAccountSize      = anchorDiscriminatorSize + 32 + 32 + 1 + 8 + 32 + 1
	licenseStatusActive     = byte(0)
	licenseStatusRevoked    = byte(1)
	licenseStatusExpired    = byte(2)
	licenseAccountName      = "License"
)

type SolanaRepository struct {
	rpcURL         string
	cluster        string
	programID      string
	programIDBytes [32]byte
	adminKeypair   []byte
	httpClient     *http.Client
}

func NewSolanaRepository(cfg config.Config) (*SolanaRepository, error) {
	if strings.TrimSpace(cfg.SolanaRPCURL) == "" {
		return nil, fmt.Errorf("SOLANA_RPC_URL is required: %w", ErrInvalidSolanaConfig)
	}

	if strings.TrimSpace(cfg.ProgramID) == "" {
		return nil, fmt.Errorf("PROGRAM_ID is required: %w", ErrInvalidSolanaConfig)
	}

	programIDBytes, err := parsePublicKey(cfg.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("PROGRAM_ID invalid: %w", err)
	}

	var adminKeypair []byte
	if strings.TrimSpace(cfg.AdminKeypairPath) != "" {
		adminKeypair, err = loadKeypair(cfg.AdminKeypairPath)
		if err != nil {
			return nil, fmt.Errorf("ADMIN_KEYPAIR_PATH invalid: %w", err)
		}
	}

	return &SolanaRepository{
		rpcURL:         cfg.SolanaRPCURL,
		cluster:        cfg.SolanaCluster,
		programID:      cfg.ProgramID,
		programIDBytes: programIDBytes,
		adminKeypair:   adminKeypair,
		httpClient:     &http.Client{Timeout: 15 * time.Second},
	}, nil
}

func (r *SolanaRepository) RegisterIssuer(ctx context.Context, params models.RegisterIssuerParams) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	return "", ErrSolanaIntegrationPending
}

func (r *SolanaRepository) IssueLicense(ctx context.Context, params models.IssueLicenseParams) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	return "", ErrSolanaIntegrationPending
}

func (r *SolanaRepository) FetchLicense(ctx context.Context, licenseHash models.LicenseHash) (*models.LicenseAccount, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	licensePDA, err := findProgramAddress([][]byte{[]byte("license"), licenseHash.Bytes[:]}, r.programIDBytes)
	if err != nil {
		return nil, err
	}

	accountData, err := r.getAccountData(ctx, encodeBase58(licensePDA[:]))
	if err != nil {
		return nil, err
	}

	return decodeLicenseAccount(accountData, licenseHash)
}

func (r *SolanaRepository) RevokeLicense(ctx context.Context, params models.RevokeLicenseParams) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	return "", ErrSolanaIntegrationPending
}

func loadKeypair(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read keypair file: %w", err)
	}

	var keyNumbers []int
	if err := json.Unmarshal(content, &keyNumbers); err != nil {
		return nil, fmt.Errorf("failed to parse keypair JSON: %w", err)
	}

	if len(keyNumbers) != 32 && len(keyNumbers) != 64 {
		return nil, fmt.Errorf("expected keypair file to contain 32 or 64 bytes, got %d", len(keyNumbers))
	}

	keypair := make([]byte, len(keyNumbers))
	for i, value := range keyNumbers {
		if value < 0 || value > 255 {
			return nil, fmt.Errorf("invalid keypair byte at index %d", i)
		}
		keypair[i] = byte(value)
	}

	return keypair, nil
}

func (r *SolanaRepository) getAccountData(ctx context.Context, account string) ([]byte, error) {
	requestBody := solanaRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getAccountInfo",
		Params: []any{
			account,
			map[string]string{
				"encoding":   "base64",
				"commitment": "confirmed",
			},
		},
	}

	rawBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, r.rpcURL, bytes.NewReader(rawBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := r.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("solana rpc unavailable: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("solana rpc returned status %d", response.StatusCode)
	}

	var rpcResponse getAccountInfoResponse
	if err := json.NewDecoder(response.Body).Decode(&rpcResponse); err != nil {
		return nil, fmt.Errorf("failed to decode solana rpc response: %w", err)
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("solana rpc error %d: %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	if rpcResponse.Result.Value == nil {
		return nil, ErrLicenseNotFound
	}

	if len(rpcResponse.Result.Value.Data) < 1 {
		return nil, fmt.Errorf("solana account data missing")
	}

	accountData, err := base64.StdEncoding.DecodeString(rpcResponse.Result.Value.Data[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode solana account data: %w", err)
	}

	return accountData, nil
}

func decodeLicenseAccount(data []byte, licenseHash models.LicenseHash) (*models.LicenseAccount, error) {
	if len(data) < licenseAccountSize {
		return nil, fmt.Errorf("license account data too short: got %d bytes", len(data))
	}

	if !bytes.Equal(data[:anchorDiscriminatorSize], anchorAccountDiscriminator(licenseAccountName)) {
		return nil, fmt.Errorf("invalid license account discriminator")
	}

	offset := anchorDiscriminatorSize
	holder := data[offset : offset+32]
	offset += 32
	issuer := data[offset : offset+32]
	offset += 32
	status := data[offset]
	offset += 1
	if status != licenseStatusActive && status != licenseStatusRevoked && status != licenseStatusExpired {
		return nil, fmt.Errorf("unknown license status discriminant %d", status)
	}

	expiry := int64(binary.LittleEndian.Uint64(data[offset : offset+8]))
	offset += 8
	assetHash := data[offset : offset+32]

	if !bytes.Equal(assetHash, licenseHash.Bytes[:]) {
		return nil, fmt.Errorf("license account hash mismatch")
	}

	expiryDate := "9999-12-31"
	if expiry > 0 {
		expiryDate = time.Unix(expiry, 0).UTC().Format("2006-01-02")
	}
	if status == licenseStatusExpired {
		expiryDate = "1970-01-01"
	}

	return &models.LicenseAccount{
		LicenseHash:  licenseHash,
		HolderWallet: encodeBase58(holder),
		IssuerWallet: encodeBase58(issuer),
		ExpiryDate:   expiryDate,
		IsRevoked:    status == licenseStatusRevoked,
	}, nil
}

func anchorAccountDiscriminator(accountName string) []byte {
	hash := sha256.Sum256([]byte("account:" + accountName))
	return hash[:anchorDiscriminatorSize]
}

func findProgramAddress(seeds [][]byte, programID [32]byte) ([32]byte, error) {
	for bump := 255; bump >= 0; bump-- {
		address, err := createProgramAddress(append(seeds, []byte{byte(bump)}), programID)
		if err == nil {
			return address, nil
		}
	}

	return [32]byte{}, fmt.Errorf("unable to find viable program address")
}

func createProgramAddress(seeds [][]byte, programID [32]byte) ([32]byte, error) {
	var buffer []byte
	for _, seed := range seeds {
		if len(seed) > 32 {
			return [32]byte{}, fmt.Errorf("pda seed cannot exceed 32 bytes")
		}
		buffer = append(buffer, seed...)
	}
	buffer = append(buffer, programID[:]...)
	buffer = append(buffer, []byte("ProgramDerivedAddress")...)

	hash := sha256.Sum256(buffer)
	if isEd25519Point(hash[:]) {
		return [32]byte{}, fmt.Errorf("derived address is on curve")
	}

	return hash, nil
}

func parsePublicKey(value string) ([32]byte, error) {
	var publicKey [32]byte

	decoded, ok := decodeBase58(value)
	if !ok || len(decoded) != 32 {
		return publicKey, fmt.Errorf("public key must be a 32-byte base58-encoded Solana address")
	}

	copy(publicKey[:], decoded)
	return publicKey, nil
}

func decodeBase58(value string) ([]byte, bool) {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	bytes := []byte{0}
	for _, r := range value {
		index := strings.IndexRune(alphabet, r)
		if index < 0 {
			return nil, false
		}

		carry := index
		for i := len(bytes) - 1; i >= 0; i-- {
			carry += int(bytes[i]) * 58
			bytes[i] = byte(carry % 256)
			carry /= 256
		}

		for carry > 0 {
			bytes = append([]byte{byte(carry % 256)}, bytes...)
			carry /= 256
		}
	}

	for _, r := range value {
		if r != '1' {
			break
		}
		bytes = append([]byte{0}, bytes...)
	}

	if len(bytes) > 1 && bytes[0] == 0 {
		bytes = bytes[1:]
	}

	return bytes, true
}

func encodeBase58(value []byte) string {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	x := new(big.Int).SetBytes(value)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)
	var encoded []byte

	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		encoded = append([]byte{alphabet[mod.Int64()]}, encoded...)
	}

	for _, b := range value {
		if b != 0 {
			break
		}
		encoded = append([]byte{alphabet[0]}, encoded...)
	}

	if len(encoded) == 0 {
		return string(alphabet[0])
	}

	return string(encoded)
}

func isEd25519Point(value []byte) bool {
	if len(value) != 32 {
		return false
	}

	encodedY := append([]byte(nil), value...)
	encodedY[31] &= 0x7f
	y := littleEndianBytesToInt(encodedY)

	p := ed25519Prime()
	if y.Cmp(p) >= 0 {
		return false
	}

	one := big.NewInt(1)
	ySquared := new(big.Int).Mul(y, y)
	ySquared.Mod(ySquared, p)

	u := new(big.Int).Sub(ySquared, one)
	u.Mod(u, p)

	d := ed25519D()
	v := new(big.Int).Mul(d, ySquared)
	v.Add(v, one)
	v.Mod(v, p)

	vInverse := new(big.Int).ModInverse(v, p)
	if vInverse == nil {
		return false
	}

	xSquared := new(big.Int).Mul(u, vInverse)
	xSquared.Mod(xSquared, p)

	return hasModSqrt(xSquared, p)
}

func hasModSqrt(value *big.Int, p *big.Int) bool {
	if value.Sign() == 0 {
		return true
	}

	exponent := new(big.Int).Add(p, big.NewInt(3))
	exponent.Div(exponent, big.NewInt(8))

	root := new(big.Int).Exp(value, exponent, p)
	check := new(big.Int).Mul(root, root)
	check.Mod(check, p)
	if check.Cmp(value) == 0 {
		return true
	}

	negValue := new(big.Int).Neg(value)
	negValue.Mod(negValue, p)
	return check.Cmp(negValue) == 0
}

func littleEndianBytesToInt(value []byte) *big.Int {
	reversed := make([]byte, len(value))
	for i := range value {
		reversed[len(value)-1-i] = value[i]
	}
	return new(big.Int).SetBytes(reversed)
}

func ed25519Prime() *big.Int {
	p := new(big.Int).Lsh(big.NewInt(1), 255)
	p.Sub(p, big.NewInt(19))
	return p
}

func ed25519D() *big.Int {
	p := ed25519Prime()
	numerator := big.NewInt(-121665)
	denominator := big.NewInt(121666)
	denominatorInverse := new(big.Int).ModInverse(denominator, p)
	d := new(big.Int).Mul(numerator, denominatorInverse)
	d.Mod(d, p)
	return d
}

type solanaRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type getAccountInfoResponse struct {
	Result getAccountInfoResult `json:"result"`
	Error  *solanaRPCError      `json:"error"`
}

type getAccountInfoResult struct {
	Value *solanaAccountValue `json:"value"`
}

type solanaAccountValue struct {
	Data []string `json:"data"`
}

type solanaRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
