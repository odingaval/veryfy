package models

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

const expiryDateLayout = "2006-01-02"

// Validate checks the issue request shape without mutating any fields.
func (r IssueLicenseRequest) Validate() error {
	var validation ValidationErrors

	validation.Required("licenseNumber", r.LicenseNumber)
	validation.Required("holderName", r.HolderName)
	validation.Required("holderWallet", r.HolderWallet)
	validation.Required("issuerWallet", r.IssuerWallet)
	validation.LicenseType("licenseType", r.LicenseType)
	validation.ExpiryDate("expiryDate", r.ExpiryDate)
	validation.Wallet("holderWallet", r.HolderWallet)
	validation.Wallet("issuerWallet", r.IssuerWallet)

	return validation.Err()
}

// Validate checks the verify request shape without mutating any fields.
func (r VerifyLicenseRequest) Validate() error {
	var validation ValidationErrors

	validation.Required("licenseNumber", r.LicenseNumber)
	validation.Required("holderName", r.HolderName)
	validation.Required("issuerWallet", r.IssuerWallet)
	validation.LicenseType("licenseType", r.LicenseType)
	validation.ExpiryDate("expiryDate", r.ExpiryDate)
	validation.Wallet("issuerWallet", r.IssuerWallet)

	return validation.Err()
}

// Validate checks the revoke request shape without mutating any fields.
func (r RevokeLicenseRequest) Validate() error {
	var validation ValidationErrors

	validation.Required("licenseHash", r.LicenseHash)
	validation.Required("issuerWallet", r.IssuerWallet)
	validation.HexHash("licenseHash", r.LicenseHash)
	validation.Wallet("issuerWallet", r.IssuerWallet)

	return validation.Err()
}

// Validate checks the issuer registration request shape without mutating any fields.
func (r RegisterIssuerRequest) Validate() error {
	var validation ValidationErrors

	validation.Required("name", r.Name)
	validation.Required("wallet", r.Wallet)
	validation.LicenseType("licenseType", r.LicenseType)
	validation.Wallet("wallet", r.Wallet)

	return validation.Err()
}

// ParseLicenseHash decodes the API hex representation into raw hash bytes.
func ParseLicenseHash(value string) (LicenseHash, error) {
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != 32 {
		return LicenseHash{}, fmt.Errorf("licenseHash must be a 64-character hex string")
	}

	var bytes [32]byte
	copy(bytes[:], decoded)

	return LicenseHash{
		Bytes: bytes,
		Hex:   strings.ToLower(value),
	}, nil
}

// IsValidLicenseType reports whether value is a clean issuer-defined type slug.
func IsValidLicenseType(value LicenseType) bool {
	text := string(value)
	if len(text) < 2 || len(text) > 64 {
		return false
	}

	for _, r := range text {
		if r >= 'A' && r <= 'Z' {
			continue
		}

		if r >= '0' && r <= '9' {
			continue
		}

		if r == '_' || r == '-' {
			continue
		}

		return false
	}

	return true
}

// IsValidExpiryDate reports whether value uses YYYY-MM-DD.
func IsValidExpiryDate(value string) bool {
	parsed, err := time.Parse(expiryDateLayout, value)
	if err != nil {
		return false
	}

	return parsed.Format(expiryDateLayout) == value
}

// ValidationErrors accumulates request validation failures.
type ValidationErrors struct {
	messages []string
}

// Required checks that a string field is not blank.
func (v *ValidationErrors) Required(field string, value string) {
	if strings.TrimSpace(value) == "" {
		v.add("%s is required", field)
	}
}

// LicenseType checks that a field contains a supported license type.
func (v *ValidationErrors) LicenseType(field string, value LicenseType) {
	if !IsValidLicenseType(value) {
		v.add("%s must be 2-64 characters using uppercase letters, numbers, underscores, or hyphens", field)
	}
}

// ExpiryDate checks that a field uses YYYY-MM-DD.
func (v *ValidationErrors) ExpiryDate(field string, value string) {
	if !IsValidExpiryDate(value) {
		v.add("%s must use YYYY-MM-DD", field)
	}
}

// Wallet checks that a field is a Solana public key encoded as base58.
func (v *ValidationErrors) Wallet(field string, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}

	decoded, ok := decodeBase58(value)
	if !ok || len(decoded) != 32 {
		v.add("%s must be a valid Solana public key", field)
	}
}

// HexHash checks that a field is a 32-byte SHA-256 hash encoded as hex.
func (v *ValidationErrors) HexHash(field string, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}

	if _, err := ParseLicenseHash(value); err != nil {
		v.add("%s must be a 64-character hex string", field)
	}
}

// Err returns an aggregated validation error, or nil when validation passed.
func (v ValidationErrors) Err() error {
	if len(v.messages) == 0 {
		return nil
	}

	return errors.New(strings.Join(v.messages, "; "))
}

func (v *ValidationErrors) add(format string, args ...any) {
	v.messages = append(v.messages, fmt.Sprintf(format, args...))
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
