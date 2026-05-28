package repositories

import (
	"context"
	"errors"

	"github.com/odingaval/veryfy/api/internal/models"
)

var (
	// ErrLicenseNotFound indicates the license hash does not exist.
	ErrLicenseNotFound = errors.New("license not found")
	// ErrIssuerNotFound indicates the issuer wallet has not been registered.
	ErrIssuerNotFound = errors.New("issuer not found")
	// ErrDuplicateLicense indicates a license hash already exists.
	ErrDuplicateLicense = errors.New("duplicate license")
	// ErrDuplicateIssuer indicates an issuer wallet has already been registered.
	ErrDuplicateIssuer = errors.New("duplicate issuer")
	// ErrUnauthorizedIssuer indicates the issuer cannot perform the requested action.
	ErrUnauthorizedIssuer = errors.New("unauthorized issuer")
)

// LicenseRepository hides the storage or chain implementation for license records.
type LicenseRepository interface {
	IssueLicense(ctx context.Context, params models.IssueLicenseParams) (string, error)
	FetchLicense(ctx context.Context, licenseHash models.LicenseHash) (*models.LicenseAccount, error)
	RevokeLicense(ctx context.Context, params models.RevokeLicenseParams) (string, error)
}

// IssuerRepository hides the storage or chain implementation for issuer records.
type IssuerRepository interface {
	RegisterIssuer(ctx context.Context, params models.RegisterIssuerParams) (string, error)
}
