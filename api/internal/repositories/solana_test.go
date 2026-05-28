package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

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

func TestParsePublicKeyRejectsInvalidValues(t *testing.T) {
	if _, err := parsePublicKey("not-a-valid-base58-key"); err == nil {
		t.Fatal("expected error for invalid public key")
	}
}

func TestSolanaRepositoryMethodsReturnIntegrationPending(t *testing.T) {
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

	if _, err := repo.FetchLicense(context.Background(), models.LicenseHash{Hex: "abc"}); !errors.Is(err, ErrSolanaIntegrationPending) {
		t.Fatalf("FetchLicense error = %v, want ErrSolanaIntegrationPending", err)
	}

	if _, err := repo.RevokeLicense(context.Background(), models.RevokeLicenseParams{LicenseHash: models.LicenseHash{Hex: "abc"}, IssuerWallet: "11111111111111111111111111111111"}); !errors.Is(err, ErrSolanaIntegrationPending) {
		t.Fatalf("RevokeLicense error = %v, want ErrSolanaIntegrationPending", err)
	}
}
