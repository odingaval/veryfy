package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/odingaval/veryfy/api/internal/models"
	"github.com/odingaval/veryfy/api/internal/repositories"
)

type stubIssuerService struct {
	response models.RegisterIssuerResponse
	err      error
}

func (s stubIssuerService) RegisterIssuer(ctx context.Context, request models.RegisterIssuerRequest) (models.RegisterIssuerResponse, error) {
	return s.response, s.err
}

type stubLicenseService struct {
	issueResponse  models.IssueLicenseResponse
	verifyResponse models.VerifyLicenseResponse
	revokeResponse models.RevokeLicenseResponse
	err            error
}

func (s stubLicenseService) IssueLicense(ctx context.Context, request models.IssueLicenseRequest) (models.IssueLicenseResponse, error) {
	return s.issueResponse, s.err
}

func (s stubLicenseService) VerifyLicense(ctx context.Context, request models.VerifyLicenseRequest) (models.VerifyLicenseResponse, error) {
	return s.verifyResponse, s.err
}

func (s stubLicenseService) RevokeLicense(ctx context.Context, request models.RevokeLicenseRequest) (models.RevokeLicenseResponse, error) {
	return s.revokeResponse, s.err
}

func TestRegisterIssuerHandler_Success(t *testing.T) {
	handler := NewRegisterIssuerHandler(stubIssuerService{
		response: models.RegisterIssuerResponse{TxSignature: "fake-register-signature"},
	})

	request := httptest.NewRequest(http.MethodPost, "/issuers/register", strings.NewReader(`{"name":"KMPDC","licenseType":"MEDICAL","wallet":"11111111111111111111111111111111"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var body struct {
		Data struct {
			TxSignature string `json:"txSignature"`
		} `json:"data"`
		Error any `json:"error"`
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Data.TxSignature != "fake-register-signature" {
		t.Fatalf("txSignature = %q, want %q", body.Data.TxSignature, "fake-register-signature")
	}

	if body.Error != nil {
		t.Fatalf("error payload = %v, want nil", body.Error)
	}
}

func TestRegisterIssuerHandler_DuplicateIssuer(t *testing.T) {
	handler := NewRegisterIssuerHandler(stubIssuerService{err: repositories.ErrDuplicateIssuer})

	request := httptest.NewRequest(http.MethodPost, "/issuers/register", strings.NewReader(`{"name":"KMPDC","licenseType":"MEDICAL","wallet":"11111111111111111111111111111111"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusConflict)
	}
}

func TestIssueLicenseHandler_InvalidJSON(t *testing.T) {
	handler := NewIssueLicenseHandler(stubLicenseService{})

	request := httptest.NewRequest(http.MethodPost, "/licenses/issue", strings.NewReader(`{"licenseNumber":"ABC"`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestVerifyLicenseHandler_Success(t *testing.T) {
	handler := NewVerifyLicenseHandler(stubLicenseService{
		verifyResponse: models.VerifyLicenseResponse{Status: models.LicenseStatusValid, LicenseHash: "abc123"},
	})

	request := httptest.NewRequest(http.MethodPost, "/licenses/verify", strings.NewReader(`{"licenseNumber":"KE/MED/12345","holderName":"John Kamau","licenseType":"MEDICAL","expiryDate":"2026-12-31","issuerWallet":"11111111111111111111111111111111"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
}

func TestRevokeLicenseHandler_Unauthorized(t *testing.T) {
	handler := NewRevokeLicenseHandler(stubLicenseService{err: repositories.ErrUnauthorizedIssuer})

	request := httptest.NewRequest(http.MethodPost, "/licenses/revoke", strings.NewReader(`{"licenseHash":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","issuerWallet":"11111111111111111111111111111111"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
}
