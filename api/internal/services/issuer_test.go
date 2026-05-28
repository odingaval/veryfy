package services

import (
	"context"
	"strings"
	"testing"

	"github.com/odingaval/veryfy/api/internal/models"
	"github.com/odingaval/veryfy/api/internal/repositories"
)

func TestIssuerServiceRegisterIssuer(t *testing.T) {
	repository := repositories.NewMemoryRepository()
	issuerService := NewIssuerService(repository)

	response, err := issuerService.RegisterIssuer(context.Background(), validIssuerRequest())
	if err != nil {
		t.Fatalf("RegisterIssuer() error = %v", err)
	}

	if !strings.HasPrefix(response.TxSignature, "fake-register-") {
		t.Fatalf("TxSignature = %q, want fake register signature", response.TxSignature)
	}
}

func TestIssuerServiceRegisterIssuerRejectsInvalidRequest(t *testing.T) {
	repository := repositories.NewMemoryRepository()
	issuerService := NewIssuerService(repository)

	_, err := issuerService.RegisterIssuer(context.Background(), validIssuerRequestWithInvalidWallet())
	if err == nil {
		t.Fatal("RegisterIssuer() error = nil, want error")
	}
}

func validIssuerRequestWithInvalidWallet() models.RegisterIssuerRequest {
	request := validIssuerRequest()
	request.Wallet = "not-a-wallet"

	return request
}
