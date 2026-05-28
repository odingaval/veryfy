package config

import (
	"reflect"
	"testing"
	"time"
)

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("ALLOWED_ORIGINS", "")
	t.Setenv("REQUEST_TIMEOUT_SECONDS", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "8080" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "8080")
	}

	if cfg.RequestTimeout != 15*time.Second {
		t.Fatalf("RequestTimeout = %s, want %s", cfg.RequestTimeout, 15*time.Second)
	}

	wantOrigins := []string{"http://localhost:5173"}
	if !reflect.DeepEqual(cfg.AllowedOrigins, wantOrigins) {
		t.Fatalf("AllowedOrigins = %v, want %v", cfg.AllowedOrigins, wantOrigins)
	}
}

func TestLoadReadsEnvironment(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:5173, http://localhost:3000")
	t.Setenv("REQUEST_TIMEOUT_SECONDS", "7")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != "9090" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "9090")
	}

	if cfg.RequestTimeout != 7*time.Second {
		t.Fatalf("RequestTimeout = %s, want %s", cfg.RequestTimeout, 7*time.Second)
	}

	wantOrigins := []string{"http://localhost:5173", "http://localhost:3000"}
	if !reflect.DeepEqual(cfg.AllowedOrigins, wantOrigins) {
		t.Fatalf("AllowedOrigins = %v, want %v", cfg.AllowedOrigins, wantOrigins)
	}
}

func TestLoadRejectsInvalidTimeout(t *testing.T) {
	t.Setenv("REQUEST_TIMEOUT_SECONDS", "nope")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want error")
	}
}

func TestLoadReadsSolanaConfiguration(t *testing.T) {
	t.Setenv("SOLANA_RPC_URL", "https://api.devnet.solana.com")
	t.Setenv("PROGRAM_ID", "11111111111111111111111111111111")
	t.Setenv("ADMIN_KEYPAIR_PATH", "/tmp/admin-keypair.json")
	t.Setenv("SOLANA_CLUSTER", "devnet")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.SolanaRPCURL != "https://api.devnet.solana.com" {
		t.Fatalf("SolanaRPCURL = %q, want %q", cfg.SolanaRPCURL, "https://api.devnet.solana.com")
	}

	if cfg.ProgramID != "11111111111111111111111111111111" {
		t.Fatalf("ProgramID = %q, want %q", cfg.ProgramID, "11111111111111111111111111111111111")
	}

	if cfg.AdminKeypairPath != "/tmp/admin-keypair.json" {
		t.Fatalf("AdminKeypairPath = %q, want %q", cfg.AdminKeypairPath, "/tmp/admin-keypair.json")
	}

	if cfg.SolanaCluster != "devnet" {
		t.Fatalf("SolanaCluster = %q, want %q", cfg.SolanaCluster, "devnet")
	}
}
