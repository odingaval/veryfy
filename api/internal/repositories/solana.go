package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	if strings.TrimSpace(cfg.AdminKeypairPath) == "" {
		return nil, fmt.Errorf("ADMIN_KEYPAIR_PATH is required: %w", ErrInvalidSolanaConfig)
	}

	programIDBytes, err := parsePublicKey(cfg.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("PROGRAM_ID invalid: %w", err)
	}

	adminKeypair, err := loadKeypair(cfg.AdminKeypairPath)
	if err != nil {
		return nil, fmt.Errorf("ADMIN_KEYPAIR_PATH invalid: %w", err)
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

	return nil, ErrSolanaIntegrationPending
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
