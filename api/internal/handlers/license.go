package handlers

import (
	"context"
	"net/http"

	"github.com/odingaval/veryfy/api/internal/httpjson"
	"github.com/odingaval/veryfy/api/internal/models"
)

type licenseOperator interface {
	IssueLicense(ctx context.Context, request models.IssueLicenseRequest) (models.IssueLicenseResponse, error)
	VerifyLicense(ctx context.Context, request models.VerifyLicenseRequest) (models.VerifyLicenseResponse, error)
	RevokeLicense(ctx context.Context, request models.RevokeLicenseRequest) (models.RevokeLicenseResponse, error)
}

type IssueLicenseHandler struct {
	service licenseOperator
}

type VerifyLicenseHandler struct {
	service licenseOperator
}

type RevokeLicenseHandler struct {
	service licenseOperator
}

func NewIssueLicenseHandler(service licenseOperator) *IssueLicenseHandler {
	return &IssueLicenseHandler{service: service}
}

func NewVerifyLicenseHandler(service licenseOperator) *VerifyLicenseHandler {
	return &VerifyLicenseHandler{service: service}
}

func NewRevokeLicenseHandler(service licenseOperator) *RevokeLicenseHandler {
	return &RevokeLicenseHandler{service: service}
}

func (h *IssueLicenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	request, err := decodeJSONBody[models.IssueLicenseRequest](r)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid JSON payload")
		return
	}

	response, err := h.service.IssueLicense(r.Context(), request)
	if err != nil {
		writeHandlerError(w, err)
		return
	}

	httpjson.WriteSuccess(w, http.StatusOK, response)
}

func (h *VerifyLicenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	request, err := decodeJSONBody[models.VerifyLicenseRequest](r)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid JSON payload")
		return
	}

	response, err := h.service.VerifyLicense(r.Context(), request)
	if err != nil {
		writeHandlerError(w, err)
		return
	}

	httpjson.WriteSuccess(w, http.StatusOK, response)
}

func (h *RevokeLicenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	request, err := decodeJSONBody[models.RevokeLicenseRequest](r)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid JSON payload")
		return
	}

	response, err := h.service.RevokeLicense(r.Context(), request)
	if err != nil {
		writeHandlerError(w, err)
		return
	}

	httpjson.WriteSuccess(w, http.StatusOK, response)
}
