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
