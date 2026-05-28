package models

import (
	"strings"
	"testing"
)

const validWallet = "11111111111111111111111111111111"

func TestIsValidLicenseType(t *testing.T) {
	tests := []struct {
		value LicenseType
		want  bool
	}{
		{"MEDICAL", true},
		{"LEGAL", true},
		{"DRIVING", true},
		{"ENGINEERING", true},
		{"PHARMACY", true},
		{"NURSING_COUNCIL", true},
		{"PILOT-LICENSE", true},
		{"A1", true},
		{"medical", false},
		{"Medical", false},
		{"MEDICAL LICENSE", false},
		{"MEDICAL/LICENSE", false},
		{"A", false},
		{LicenseType(strings.Repeat("A", 65)), false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsValidLicenseType(tt.value); got != tt.want {
			t.Fatalf("IsValidLicenseType(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestIsValidExpiryDate(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"2026-12-31", true},
		{"2026-02-29", false},
		{"31-12-2026", false},
		{"2026-12-31T00:00:00Z", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsValidExpiryDate(tt.value); got != tt.want {
			t.Fatalf("IsValidExpiryDate(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestParseLicenseHash(t *testing.T) {
	value := strings.Repeat("a", 64)

	hash, err := ParseLicenseHash(value)
	if err != nil {
		t.Fatalf("ParseLicenseHash() error = %v", err)
	}

	if hash.Hex != value {
		t.Fatalf("Hex = %q, want %q", hash.Hex, value)
	}

	if len(hash.Bytes) != 32 {
		t.Fatalf("Bytes length = %d, want %d", len(hash.Bytes), 32)
	}
}

func TestParseLicenseHashRejectsInvalidHex(t *testing.T) {
	if _, err := ParseLicenseHash("not-a-hash"); err == nil {
		t.Fatal("ParseLicenseHash() error = nil, want error")
	}
}

func TestIssueLicenseRequestValidate(t *testing.T) {
	request := IssueLicenseRequest{
		LicenseNumber: "KE/MED/12345",
		HolderName:    "John Kamau",
		HolderWallet:  validWallet,
		LicenseType:   "PHARMACY",
		ExpiryDate:    "2026-12-31",
		IssuerWallet:  validWallet,
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestVerifyLicenseRequestValidateRejectsInvalidShape(t *testing.T) {
	request := VerifyLicenseRequest{
		LicenseNumber: "",
		HolderName:    "John Kamau",
		LicenseType:   "accounting",
		ExpiryDate:    "12/31/2026",
		IssuerWallet:  "not-a-wallet",
	}

	err := request.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want error")
	}

	message := err.Error()
	for _, want := range []string{
		"licenseNumber is required",
		"licenseType must be 2-64 characters using uppercase letters, numbers, underscores, or hyphens",
		"expiryDate must use YYYY-MM-DD",
		"issuerWallet must be a valid Solana public key",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("Validate() error = %q, want to contain %q", message, want)
		}
	}
}

func TestRevokeLicenseRequestValidate(t *testing.T) {
	request := RevokeLicenseRequest{
		LicenseHash:  strings.Repeat("f", 64),
		IssuerWallet: validWallet,
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestRegisterIssuerRequestValidate(t *testing.T) {
	request := RegisterIssuerRequest{
		Name:        "KMPDC",
		LicenseType: "MEDICAL_BOARD",
		Wallet:      validWallet,
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
