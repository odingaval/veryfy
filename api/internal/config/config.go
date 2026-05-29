package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPort           = "8080"
	defaultRequestTimeout = 15 * time.Second
)

// Config holds runtime settings for the API server.
type Config struct {
	Port             string
	SolanaRPCURL     string
	ProgramID        string
	AdminKeypairPath string
	SolanaCluster    string
	AllowedOrigins   []string
	RequestTimeout   time.Duration
}

// Load reads configuration from the process environment.
func Load() (Config, error) {
	cfg := Config{
		Port:             getEnv("PORT", defaultPort),
		SolanaRPCURL:     os.Getenv("SOLANA_RPC_URL"),
		ProgramID:        os.Getenv("PROGRAM_ID"),
		AdminKeypairPath: os.Getenv("ADMIN_KEYPAIR_PATH"),
		SolanaCluster:    os.Getenv("SOLANA_CLUSTER"),
		AllowedOrigins:   splitCSV(getEnv("ALLOWED_ORIGINS", "http://localhost:5173")),
		RequestTimeout:   defaultRequestTimeout,
	}

	if rawTimeout := os.Getenv("REQUEST_TIMEOUT_SECONDS"); rawTimeout != "" {
		seconds, err := strconv.Atoi(rawTimeout)
		if err != nil || seconds <= 0 {
			return Config{}, fmt.Errorf("REQUEST_TIMEOUT_SECONDS must be a positive integer")
		}

		cfg.RequestTimeout = time.Duration(seconds) * time.Second
	}

	return cfg, nil
}

// Addr returns the TCP address used by http.Server.
func (c Config) Addr() string {
	return ":" + c.Port
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
