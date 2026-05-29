package services

import (
	"context"
	"errors"
	"time"

	"github.com/odingaval/veryfy/api/internal/models"
	"github.com/odingaval/veryfy/api/internal/repositories"
)

const serviceDateLayout = "2006-01-02"

// LicenseService coordinates license issuing, verification, and revocation.
type LicenseService struct {
	hashService *HashService
	repository  repositories.LicenseRepository
	now         func() time.Time
}

// NewLicenseService creates a license service.
func NewLicenseService(hashService *HashService, repository repositories.LicenseRepository) *LicenseService {
	return &LicenseService{
		hashService: hashService,
		repository:  repository,
		now:         time.Now,
	}
}

// IssueLicense validates, hashes, and stores a new license.
func (s *LicenseService) IssueLicense(ctx context.Context, request models.IssueLicenseRequest) (models.IssueLicenseResponse, error) {
	if err := request.Validate(); err != nil {
		return models.IssueLicenseResponse{}, err
	}

	licenseHash := s.hashService.GenerateLicenseHash(
		request.LicenseNumber,
		request.HolderName,
		request.IssuerWallet,
		request.LicenseType,
		request.ExpiryDate,
	)

	txSignature, err := s.repository.IssueLicense(ctx, models.IssueLicenseParams{
		LicenseHash:   licenseHash,
		LicenseNumber: request.LicenseNumber,
		HolderName:    request.HolderName,
		HolderWallet:  request.HolderWallet,
		IssuerWallet:  request.IssuerWallet,
		LicenseType:   request.LicenseType,
		ExpiryDate:    request.ExpiryDate,
	})
	if err != nil {
		return models.IssueLicenseResponse{}, err
	}

	return models.IssueLicenseResponse{
		LicenseHash: licenseHash.Hex,
		TxSignature: txSignature,
	}, nil
}

// VerifyLicense validates request details and returns the derived verification status.
func (s *LicenseService) VerifyLicense(ctx context.Context, request models.VerifyLicenseRequest) (models.VerifyLicenseResponse, error) {
	if err := request.Validate(); err != nil {
		return models.VerifyLicenseResponse{}, err
	}

	licenseHash := s.hashService.GenerateLicenseHash(
		request.LicenseNumber,
		request.HolderName,
		request.IssuerWallet,
		request.LicenseType,
		request.ExpiryDate,
	)

	license, err := s.repository.FetchLicense(ctx, licenseHash)
	if errors.Is(err, repositories.ErrLicenseNotFound) {
		return models.VerifyLicenseResponse{
			Status:      models.LicenseStatusInvalid,
			LicenseHash: licenseHash.Hex,
		}, nil
	}
	if err != nil {
		return models.VerifyLicenseResponse{}, err
	}

	return models.VerifyLicenseResponse{
		Status:      s.statusForLicense(*license),
		LicenseHash: licenseHash.Hex,
	}, nil
}

// RevokeLicense validates and revokes a license by hash.
func (s *LicenseService) RevokeLicense(ctx context.Context, request models.RevokeLicenseRequest) (models.RevokeLicenseResponse, error) {
	if err := request.Validate(); err != nil {
		return models.RevokeLicenseResponse{}, err
	}

	licenseHash, err := models.ParseLicenseHash(request.LicenseHash)
	if err != nil {
		return models.RevokeLicenseResponse{}, err
	}

	txSignature, err := s.repository.RevokeLicense(ctx, models.RevokeLicenseParams{
		LicenseHash:  licenseHash,
		IssuerWallet: request.IssuerWallet,
	})
	if err != nil {
		return models.RevokeLicenseResponse{}, err
	}

	return models.RevokeLicenseResponse{
		TxSignature: txSignature,
	}, nil
}

func (s *LicenseService) statusForLicense(license models.LicenseAccount) models.LicenseStatus {
	if license.IsRevoked {
		return models.LicenseStatusRevoked
	}

	expiryDate, err := time.Parse(serviceDateLayout, license.ExpiryDate)
	if err != nil {
		return models.LicenseStatusInvalid
	}

	today, err := time.Parse(serviceDateLayout, s.now().Format(serviceDateLayout))
	if err != nil {
		return models.LicenseStatusInvalid
	}

	if expiryDate.Before(today) {
		return models.LicenseStatusExpired
	}

	return models.LicenseStatusValid
}
