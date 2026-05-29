package repositories

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/odingaval/veryfy/api/internal/models"
)

// MemoryRepository is an in-memory implementation for tests and local development.
type MemoryRepository struct {
	mu       sync.RWMutex
	licenses map[string]models.LicenseAccount
	issuers  map[string]models.IssuerAccount
	now      func() time.Time
}

// NewMemoryRepository creates an empty in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		licenses: make(map[string]models.LicenseAccount),
		issuers:  make(map[string]models.IssuerAccount),
		now:      time.Now,
	}
}

// RegisterIssuer registers an issuer wallet and license type.
func (r *MemoryRepository) RegisterIssuer(ctx context.Context, params models.RegisterIssuerParams) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.issuers[params.Wallet]; exists {
		return "", ErrDuplicateIssuer
	}

	r.issuers[params.Wallet] = models.IssuerAccount{
		Name:         params.Name,
		Wallet:       params.Wallet,
		LicenseType:  params.LicenseType,
		RegisteredAt: r.now(),
	}

	return fakeSignature("register", params.Wallet), nil
}

// IssueLicense stores a license account after checking issuer registration.
func (r *MemoryRepository) IssueLicense(ctx context.Context, params models.IssueLicenseParams) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	issuer, exists := r.issuers[params.IssuerWallet]
	if !exists {
		return "", ErrIssuerNotFound
	}

	if issuer.LicenseType != params.LicenseType {
		return "", ErrUnauthorizedIssuer
	}

	if _, exists := r.licenses[params.LicenseHash.Hex]; exists {
		return "", ErrDuplicateLicense
	}

	r.licenses[params.LicenseHash.Hex] = models.LicenseAccount{
		LicenseHash:   params.LicenseHash,
		LicenseNumber: params.LicenseNumber,
		HolderName:    params.HolderName,
		HolderWallet:  params.HolderWallet,
		IssuerWallet:  params.IssuerWallet,
		LicenseType:   params.LicenseType,
		ExpiryDate:    params.ExpiryDate,
		IsRevoked:     false,
		IssuedAt:      r.now(),
	}

	return fakeSignature("issue", params.LicenseHash.Hex), nil
}

// FetchLicense returns a license account by hash.
func (r *MemoryRepository) FetchLicense(ctx context.Context, licenseHash models.LicenseHash) (*models.LicenseAccount, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	license, exists := r.licenses[licenseHash.Hex]
	if !exists {
		return nil, ErrLicenseNotFound
	}

	return &license, nil
}

// RevokeLicense marks a license as revoked when called by the original issuer.
func (r *MemoryRepository) RevokeLicense(ctx context.Context, params models.RevokeLicenseParams) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	license, exists := r.licenses[params.LicenseHash.Hex]
	if !exists {
		return "", ErrLicenseNotFound
	}

	if license.IssuerWallet != params.IssuerWallet {
		return "", ErrUnauthorizedIssuer
	}

	license.IsRevoked = true
	r.licenses[params.LicenseHash.Hex] = license

	return fakeSignature("revoke", params.LicenseHash.Hex), nil
}

func fakeSignature(action string, value string) string {
	if len(value) > 12 {
		value = value[:12]
	}

	return fmt.Sprintf("fake-%s-%s", action, value)
}
