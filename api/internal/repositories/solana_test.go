package repositories

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/odingaval/veryfy/api/internal/config"
	"github.com/odingaval/veryfy/api/internal/models"
)

func TestNewSolanaRepositoryValidatesConfig(t *testing.T) {
	cfg := config.Config{}

	if _, err := NewSolanaRepository(cfg); err == nil {
		t.Fatal("expected error for missing Solana configuration")
	}
}

func TestNewSolanaRepositoryLoadsKeypairFile(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "admin-keypair.json")

	keypair := make([]int, 64)
	for i := range keypair {
		keypair[i] = i
	}

	raw, err := json.Marshal(keypair)
	if err != nil {
		t.Fatalf("marshal keypair: %v", err)
	}

	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatalf("write keypair file: %v", err)
	}

	cfg := config.Config{
		SolanaRPCURL:     "https://api.devnet.solana.com",
		ProgramID:        "11111111111111111111111111111111",
		AdminKeypairPath: path,
		SolanaCluster:    "devnet",
	}

	repo, err := NewSolanaRepository(cfg)
	if err != nil {
		t.Fatalf("NewSolanaRepository() error = %v", err)
	}

	if repo.rpcURL != cfg.SolanaRPCURL {
		t.Fatalf("rpcURL = %q, want %q", repo.rpcURL, cfg.SolanaRPCURL)
	}

	if repo.programID != cfg.ProgramID {
		t.Fatalf("programID = %q, want %q", repo.programID, cfg.ProgramID)
	}
}

func TestNewSolanaRepositoryAllowsReadOnlyConfigWithoutKeypair(t *testing.T) {
	cfg := config.Config{
		SolanaRPCURL:  "https://api.devnet.solana.com",
		ProgramID:     "11111111111111111111111111111111",
		SolanaCluster: "devnet",
	}

	repo, err := NewSolanaRepository(cfg)
	if err != nil {
		t.Fatalf("NewSolanaRepository() error = %v", err)
	}

	if len(repo.adminKeypair) != 0 {
		t.Fatal("adminKeypair should be empty for read-only config")
	}
}

func TestParsePublicKeyRejectsInvalidValues(t *testing.T) {
	if _, err := parsePublicKey("not-a-valid-base58-key"); err == nil {
		t.Fatal("expected error for invalid public key")
	}
}

func TestSolanaRepositoryWriteMethodsReturnIntegrationPending(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "admin-keypair.json")

	keypair := make([]int, 64)
	for i := range keypair {
		keypair[i] = i
	}

	raw, err := json.Marshal(keypair)
	if err != nil {
		t.Fatalf("marshal keypair: %v", err)
	}

	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatalf("write keypair file: %v", err)
	}

	cfg := config.Config{
		SolanaRPCURL:     "https://api.devnet.solana.com",
		ProgramID:        "11111111111111111111111111111111",
		AdminKeypairPath: path,
		SolanaCluster:    "devnet",
	}

	repo, err := NewSolanaRepository(cfg)
	if err != nil {
		t.Fatalf("NewSolanaRepository() error = %v", err)
	}

	if _, err := repo.RegisterIssuer(context.Background(), models.RegisterIssuerParams{Name: "KMPDC", Wallet: "11111111111111111111111111111111", LicenseType: "MEDICAL"}); !errors.Is(err, ErrSolanaIntegrationPending) {
		t.Fatalf("RegisterIssuer error = %v, want ErrSolanaIntegrationPending", err)
	}

	if _, err := repo.IssueLicense(context.Background(), models.IssueLicenseParams{LicenseHash: models.LicenseHash{Hex: "abc"}, LicenseNumber: "A", HolderName: "B", HolderWallet: "11111111111111111111111111111111", IssuerWallet: "11111111111111111111111111111111", LicenseType: "MEDICAL", ExpiryDate: "2026-12-31"}); !errors.Is(err, ErrSolanaIntegrationPending) {
		t.Fatalf("IssueLicense error = %v, want ErrSolanaIntegrationPending", err)
	}

	if _, err := repo.RevokeLicense(context.Background(), models.RevokeLicenseParams{LicenseHash: models.LicenseHash{Hex: "abc"}, IssuerWallet: "11111111111111111111111111111111"}); !errors.Is(err, ErrSolanaIntegrationPending) {
		t.Fatalf("RevokeLicense error = %v, want ErrSolanaIntegrationPending", err)
	}
}

func TestFindProgramAddressIsStableAndOffCurve(t *testing.T) {
	programID, err := parsePublicKey("11111111111111111111111111111111")
	if err != nil {
		t.Fatalf("parse public key: %v", err)
	}

	licenseHash := models.LicenseHash{}
	for i := range licenseHash.Bytes {
		licenseHash.Bytes[i] = byte(i)
	}

	first, err := findProgramAddress([][]byte{[]byte("license"), licenseHash.Bytes[:]}, programID)
	if err != nil {
		t.Fatalf("findProgramAddress() error = %v", err)
	}

	second, err := findProgramAddress([][]byte{[]byte("license"), licenseHash.Bytes[:]}, programID)
	if err != nil {
		t.Fatalf("findProgramAddress() error = %v", err)
	}

	if first != second {
		t.Fatal("findProgramAddress() was not deterministic")
	}

	if isEd25519Point(first[:]) {
		t.Fatal("derived PDA is on curve")
	}
}

func TestDecodeLicenseAccount(t *testing.T) {
	licenseHash := testLicenseHash()
	expiry := time.Date(2027, time.December, 31, 0, 0, 0, 0, time.UTC).Unix()
	data := testLicenseAccountData(licenseHash, licenseStatusActive, expiry)

	license, err := decodeLicenseAccount(data, licenseHash)
	if err != nil {
		t.Fatalf("decodeLicenseAccount() error = %v", err)
	}

	if license.LicenseHash != licenseHash {
		t.Fatal("decoded license hash mismatch")
	}

	if license.ExpiryDate != "2027-12-31" {
		t.Fatalf("ExpiryDate = %q, want 2027-12-31", license.ExpiryDate)
	}

	if license.IsRevoked {
		t.Fatal("IsRevoked = true, want false")
	}

	if license.HolderWallet == "" || license.IssuerWallet == "" {
		t.Fatal("decoded wallets should not be empty")
	}
}

func TestDecodeLicenseAccountRevoked(t *testing.T) {
	licenseHash := testLicenseHash()
	data := testLicenseAccountData(licenseHash, licenseStatusRevoked, 0)

	license, err := decodeLicenseAccount(data, licenseHash)
	if err != nil {
		t.Fatalf("decodeLicenseAccount() error = %v", err)
	}

	if !license.IsRevoked {
		t.Fatal("IsRevoked = false, want true")
	}

	if license.ExpiryDate != "9999-12-31" {
		t.Fatalf("ExpiryDate = %q, want non-expiring sentinel", license.ExpiryDate)
	}
}

func TestDecodeLicenseAccountRejectsHashMismatch(t *testing.T) {
	licenseHash := testLicenseHash()
	data := testLicenseAccountData(licenseHash, licenseStatusActive, 0)
	otherHash := licenseHash
	otherHash.Bytes[0] = 99

	if _, err := decodeLicenseAccount(data, otherHash); err == nil {
		t.Fatal("expected hash mismatch error")
	}
}

func TestSolanaRepositoryFetchLicense(t *testing.T) {
	licenseHash := testLicenseHash()
	data := testLicenseAccountData(licenseHash, licenseStatusActive, time.Date(2028, time.January, 15, 0, 0, 0, 0, time.UTC).Unix())

	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var request solanaRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		if request.Method != "getAccountInfo" {
			t.Fatalf("Method = %q, want getAccountInfo", request.Method)
		}

		response := map[string]any{
			"jsonrpc": "2.0",
			"result": map[string]any{
				"value": map[string]any{
					"data": []string{base64.StdEncoding.EncodeToString(data), "base64"},
				},
			},
			"id": 1,
		}

		raw, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("marshal response: %v", err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(string(raw))),
			Header:     make(http.Header),
		}, nil
	})}

	repo := &SolanaRepository{
		rpcURL:         "http://solana.test",
		programIDBytes: mustParsePublicKey(t, "11111111111111111111111111111111"),
		httpClient:     client,
	}

	license, err := repo.FetchLicense(context.Background(), licenseHash)
	if err != nil {
		t.Fatalf("FetchLicense() error = %v", err)
	}

	if license.ExpiryDate != "2028-01-15" {
		t.Fatalf("ExpiryDate = %q, want 2028-01-15", license.ExpiryDate)
	}
}

func TestSolanaRepositoryFetchLicenseNotFound(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		response := map[string]any{
			"jsonrpc": "2.0",
			"result": map[string]any{
				"value": nil,
			},
			"id": 1,
		}

		raw, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("marshal response: %v", err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(string(raw))),
			Header:     make(http.Header),
		}, nil
	})}

	repo := &SolanaRepository{
		rpcURL:         "http://solana.test",
		programIDBytes: mustParsePublicKey(t, "11111111111111111111111111111111"),
		httpClient:     client,
	}

	if _, err := repo.FetchLicense(context.Background(), testLicenseHash()); !errors.Is(err, ErrLicenseNotFound) {
		t.Fatalf("FetchLicense() error = %v, want ErrLicenseNotFound", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func testLicenseHash() models.LicenseHash {
	var hash models.LicenseHash
	for i := range hash.Bytes {
		hash.Bytes[i] = byte(i + 1)
	}
	hash.Hex = "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	return hash
}

func testLicenseAccountData(licenseHash models.LicenseHash, status byte, expiry int64) []byte {
	data := make([]byte, licenseAccountSize)
	copy(data[:anchorDiscriminatorSize], anchorAccountDiscriminator(licenseAccountName))

	offset := anchorDiscriminatorSize

	for i := 0; i < 32; i++ {
		data[offset+i] = byte(i + 20)
	}
	offset += 32

	for i := 0; i < 32; i++ {
		data[offset+i] = byte(i + 90)
	}
	offset += 32

	data[offset] = status
	offset += 1

	binary.LittleEndian.PutUint64(data[offset:offset+8], uint64(expiry))
	offset += 8

	copy(data[offset:offset+32], licenseHash.Bytes[:])
	return data
}

func mustParsePublicKey(t *testing.T, value string) [32]byte {
	t.Helper()

	publicKey, err := parsePublicKey(value)
	if err != nil {
		t.Fatalf("parse public key: %v", err)
	}

	return publicKey
}
