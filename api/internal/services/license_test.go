package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/odingaval/veryfy/api/internal/models"
	"github.com/odingaval/veryfy/api/internal/repositories"
)

const (
	serviceTestWallet     = "11111111111111111111111111111111"
	serviceTestHolderName = "John Kamau"
	serviceTestLicenseNo  = "KE/MED/12345"
)

func TestLicenseServiceIssueLicense(t *testing.T) {
	licenseService, _ := newTestServices(t)

	response, err := licenseService.IssueLicense(context.Background(), validIssueRequest("2026-12-31"))
	if err != nil {
		t.Fatalf("IssueLicense() error = %v", err)
	}

	if response.LicenseHash == "" {
		t.Fatal("LicenseHash is empty")
	}

	if !strings.HasPrefix(response.TxSignature, "fake-issue-") {
		t.Fatalf("TxSignature = %q, want fake issue signature", response.TxSignature)
	}
}

func TestLicenseServiceVerifyLicenseValid(t *testing.T) {
	licenseService, _ := newTestServices(t)

	if _, err := licenseService.IssueLicense(context.Background(), validIssueRequest("2026-12-31")); err != nil {
		t.Fatalf("IssueLicense() error = %v", err)
	}

	response, err := licenseService.VerifyLicense(context.Background(), validVerifyRequest("2026-12-31"))
	if err != nil {
		t.Fatalf("VerifyLicense() error = %v", err)
	}

	if response.Status != models.LicenseStatusValid {
		t.Fatalf("Status = %q, want %q", response.Status, models.LicenseStatusValid)
	}
}

func TestLicenseServiceVerifyLicenseInvalid(t *testing.T) {
	licenseService, _ := newTestServices(t)

	response, err := licenseService.VerifyLicense(context.Background(), validVerifyRequest("2026-12-31"))
	if err != nil {
		t.Fatalf("VerifyLicense() error = %v", err)
	}

	if response.Status != models.LicenseStatusInvalid {
		t.Fatalf("Status = %q, want %q", response.Status, models.LicenseStatusInvalid)
	}

	if response.LicenseHash == "" {
		t.Fatal("LicenseHash is empty")
	}
}

func TestLicenseServiceVerifyLicenseRevoked(t *testing.T) {
	licenseService, _ := newTestServices(t)

	issueResponse, err := licenseService.IssueLicense(context.Background(), validIssueRequest("2026-12-31"))
	if err != nil {
		t.Fatalf("IssueLicense() error = %v", err)
	}

	if _, err := licenseService.RevokeLicense(context.Background(), models.RevokeLicenseRequest{
		LicenseHash:  issueResponse.LicenseHash,
		IssuerWallet: serviceTestWallet,
	}); err != nil {
		t.Fatalf("RevokeLicense() error = %v", err)
	}

	response, err := licenseService.VerifyLicense(context.Background(), validVerifyRequest("2026-12-31"))
	if err != nil {
		t.Fatalf("VerifyLicense() error = %v", err)
	}

	if response.Status != models.LicenseStatusRevoked {
		t.Fatalf("Status = %q, want %q", response.Status, models.LicenseStatusRevoked)
	}
}

func TestLicenseServiceVerifyLicenseExpired(t *testing.T) {
	licenseService, _ := newTestServices(t)

	if _, err := licenseService.IssueLicense(context.Background(), validIssueRequest("2025-12-31")); err != nil {
		t.Fatalf("IssueLicense() error = %v", err)
	}

	response, err := licenseService.VerifyLicense(context.Background(), validVerifyRequest("2025-12-31"))
	if err != nil {
		t.Fatalf("VerifyLicense() error = %v", err)
	}

	if response.Status != models.LicenseStatusExpired {
		t.Fatalf("Status = %q, want %q", response.Status, models.LicenseStatusExpired)
	}
}

func TestLicenseServiceRevokeLicense(t *testing.T) {
	licenseService, _ := newTestServices(t)

	issueResponse, err := licenseService.IssueLicense(context.Background(), validIssueRequest("2026-12-31"))
	if err != nil {
		t.Fatalf("IssueLicense() error = %v", err)
	}

	response, err := licenseService.RevokeLicense(context.Background(), models.RevokeLicenseRequest{
		LicenseHash:  issueResponse.LicenseHash,
		IssuerWallet: serviceTestWallet,
	})
	if err != nil {
		t.Fatalf("RevokeLicense() error = %v", err)
	}

	if !strings.HasPrefix(response.TxSignature, "fake-revoke-") {
		t.Fatalf("TxSignature = %q, want fake revoke signature", response.TxSignature)
	}
}

func TestLicenseServiceIssueLicenseRequiresRegisteredIssuer(t *testing.T) {
	repository := repositories.NewMemoryRepository()
	licenseService := NewLicenseService(NewHashService(), repository)

	_, err := licenseService.IssueLicense(context.Background(), validIssueRequest("2026-12-31"))
	if err == nil {
		t.Fatal("IssueLicense() error = nil, want error")
	}
}

func newTestServices(t *testing.T) (*LicenseService, *IssuerService) {
	t.Helper()

	repository := repositories.NewMemoryRepository()
	hashService := NewHashService()
	licenseService := NewLicenseService(hashService, repository)
	licenseService.now = func() time.Time {
		return time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	}

	issuerService := NewIssuerService(repository)
	if _, err := issuerService.RegisterIssuer(context.Background(), validIssuerRequest()); err != nil {
		t.Fatalf("RegisterIssuer() error = %v", err)
	}

	return licenseService, issuerService
}

func validIssueRequest(expiryDate string) models.IssueLicenseRequest {
	return models.IssueLicenseRequest{
		LicenseNumber: serviceTestLicenseNo,
		HolderName:    serviceTestHolderName,
		HolderWallet:  serviceTestWallet,
		LicenseType:   "MEDICAL",
		ExpiryDate:    expiryDate,
		IssuerWallet:  serviceTestWallet,
	}
}

func validVerifyRequest(expiryDate string) models.VerifyLicenseRequest {
	return models.VerifyLicenseRequest{
		LicenseNumber: serviceTestLicenseNo,
		HolderName:    serviceTestHolderName,
		LicenseType:   "MEDICAL",
		ExpiryDate:    expiryDate,
		IssuerWallet:  serviceTestWallet,
	}
}

func validIssuerRequest() models.RegisterIssuerRequest {
	return models.RegisterIssuerRequest{
		Name:        "KMPDC",
		LicenseType: "MEDICAL",
		Wallet:      serviceTestWallet,
	}
}
