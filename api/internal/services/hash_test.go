package services

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/odingaval/veryfy/api/internal/models"
)

const testIssuerWallet = "11111111111111111111111111111111"

func TestHashServiceInputStringUsesSharedContract(t *testing.T) {
	service := NewHashService()

	input := service.InputString(
		"KE/MED/12345",
		"John Kamau",
		testIssuerWallet,
		models.LicenseType("MEDICAL"),
		"2026-12-31",
	)

	want := "KE/MED/12345|John Kamau|11111111111111111111111111111111|MEDICAL|2026-12-31"
	if input != want {
		t.Fatalf("InputString() = %q, want %q", input, want)
	}
}

func TestHashServiceGenerateLicenseHashIsDeterministic(t *testing.T) {
	service := NewHashService()

	first := service.GenerateLicenseHash("KE/MED/12345", "John Kamau", testIssuerWallet, models.LicenseType("MEDICAL"), "2026-12-31")
	second := service.GenerateLicenseHash("KE/MED/12345", "John Kamau", testIssuerWallet, models.LicenseType("MEDICAL"), "2026-12-31")

	if first != second {
		t.Fatalf("hashes differ: %v != %v", first, second)
	}

	sum := sha256.Sum256([]byte("KE/MED/12345|John Kamau|11111111111111111111111111111111|MEDICAL|2026-12-31"))
	wantHex := hex.EncodeToString(sum[:])
	if first.Hex != wantHex {
		t.Fatalf("hash hex = %q, want %q", first.Hex, wantHex)
	}
}

func TestHashServiceGenerateLicenseHashChangesWhenFieldsChange(t *testing.T) {
	service := NewHashService()

	base := service.GenerateLicenseHash("KE/MED/12345", "John Kamau", testIssuerWallet, models.LicenseType("MEDICAL"), "2026-12-31")

	tests := []struct {
		name          string
		licenseNumber string
		holderName    string
		issuerWallet  string
		licenseType   models.LicenseType
		expiryDate    string
	}{
		{
			name:          "license number",
			licenseNumber: "KE/MED/54321",
			holderName:    "John Kamau",
			issuerWallet:  testIssuerWallet,
			licenseType:   models.LicenseType("MEDICAL"),
			expiryDate:    "2026-12-31",
		},
		{
			name:          "holder name",
			licenseNumber: "KE/MED/12345",
			holderName:    "Jon Kamau",
			issuerWallet:  testIssuerWallet,
			licenseType:   models.LicenseType("MEDICAL"),
			expiryDate:    "2026-12-31",
		},
		{
			name:          "issuer wallet",
			licenseNumber: "KE/MED/12345",
			holderName:    "John Kamau",
			issuerWallet:  "11111111111111111111111111111112",
			licenseType:   models.LicenseType("MEDICAL"),
			expiryDate:    "2026-12-31",
		},
		{
			name:          "license type",
			licenseNumber: "KE/MED/12345",
			holderName:    "John Kamau",
			issuerWallet:  testIssuerWallet,
			licenseType:   models.LicenseType("LEGAL"),
			expiryDate:    "2026-12-31",
		},
		{
			name:          "expiry date",
			licenseNumber: "KE/MED/12345",
			holderName:    "John Kamau",
			issuerWallet:  testIssuerWallet,
			licenseType:   models.LicenseType("MEDICAL"),
			expiryDate:    "2027-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := service.GenerateLicenseHash(tt.licenseNumber, tt.holderName, tt.issuerWallet, tt.licenseType, tt.expiryDate)
			if hash == base {
				t.Fatalf("hash did not change when %s changed", tt.name)
			}
		})
	}
}
