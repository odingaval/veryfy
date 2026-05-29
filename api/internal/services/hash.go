package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/odingaval/veryfy/api/internal/models"
)

// HashService generates deterministic license hashes.
type HashService struct{}

// NewHashService creates a hash service.
func NewHashService() *HashService {
	return &HashService{}
}

// GenerateLicenseHash hashes the shared license identity fields.
func (s *HashService) GenerateLicenseHash(
	licenseNumber string,
	holderName string,
	issuerWallet string,
	licenseType models.LicenseType,
	expiryDate string,
) models.LicenseHash {
	input := s.InputString(licenseNumber, holderName, issuerWallet, licenseType, expiryDate)
	sum := sha256.Sum256([]byte(input))

	return models.LicenseHash{
		Bytes: sum,
		Hex:   hex.EncodeToString(sum[:]),
	}
}

// InputString returns the exact shared hash input string.
func (s *HashService) InputString(
	licenseNumber string,
	holderName string,
	issuerWallet string,
	licenseType models.LicenseType,
	expiryDate string,
) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", licenseNumber, holderName, issuerWallet, licenseType, expiryDate)
}
