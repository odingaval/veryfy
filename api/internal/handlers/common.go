package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/odingaval/veryfy/api/internal/httpjson"
	"github.com/odingaval/veryfy/api/internal/repositories"
)

func decodeJSONBody[T any](r *http.Request) (T, error) {
	var value T
	if r.Body == nil {
		return value, errors.New("empty request body")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		return value, err
	}

	return value, nil
}

func writeHandlerError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, repositories.ErrDuplicateIssuer):
		httpjson.WriteError(w, http.StatusConflict, "DUPLICATE_ISSUER", "issuer already exists")
	case errors.Is(err, repositories.ErrDuplicateLicense):
		httpjson.WriteError(w, http.StatusConflict, "DUPLICATE_LICENSE", "license already exists")
	case errors.Is(err, repositories.ErrUnauthorizedIssuer):
		httpjson.WriteError(w, http.StatusForbidden, "UNAUTHORIZED", "issuer is not authorized")
	case errors.Is(err, repositories.ErrIssuerNotFound):
		httpjson.WriteError(w, http.StatusNotFound, "NOT_FOUND", "issuer not found")
	case errors.Is(err, repositories.ErrLicenseNotFound):
		httpjson.WriteError(w, http.StatusNotFound, "NOT_FOUND", "license not found")
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		httpjson.WriteError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "request cancelled or timed out")
	case isValidationError(err):
		httpjson.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
	default:
		httpjson.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

func isValidationError(err error) bool {
	if err == nil {
		return false
	}

	message := err.Error()
	return strings.Contains(message, "is required") || strings.Contains(message, "must") || strings.Contains(message, "invalid")
}
