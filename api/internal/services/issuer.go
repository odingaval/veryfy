package services

import (
	"context"

	"github.com/odingaval/veryfy/api/internal/models"
	"github.com/odingaval/veryfy/api/internal/repositories"
)

// IssuerService coordinates issuer registration.
type IssuerService struct {
	repository repositories.IssuerRepository
}

// NewIssuerService creates an issuer service.
func NewIssuerService(repository repositories.IssuerRepository) *IssuerService {
	return &IssuerService{
		repository: repository,
	}
}

// RegisterIssuer validates and registers an issuer.
func (s *IssuerService) RegisterIssuer(ctx context.Context, request models.RegisterIssuerRequest) (models.RegisterIssuerResponse, error) {
	if err := request.Validate(); err != nil {
		return models.RegisterIssuerResponse{}, err
	}

	txSignature, err := s.repository.RegisterIssuer(ctx, models.RegisterIssuerParams{
		Name:        request.Name,
		Wallet:      request.Wallet,
		LicenseType: request.LicenseType,
	})
	if err != nil {
		return models.RegisterIssuerResponse{}, err
	}

	return models.RegisterIssuerResponse{
		TxSignature: txSignature,
	}, nil
}
