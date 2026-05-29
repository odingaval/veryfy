package handlers

import (
	"context"
	"net/http"

	"github.com/odingaval/veryfy/api/internal/httpjson"
	"github.com/odingaval/veryfy/api/internal/models"
)

type issuerRegistrar interface {
	RegisterIssuer(ctx context.Context, request models.RegisterIssuerRequest) (models.RegisterIssuerResponse, error)
}

type RegisterIssuerHandler struct {
	service issuerRegistrar
}

func NewRegisterIssuerHandler(service issuerRegistrar) *RegisterIssuerHandler {
	return &RegisterIssuerHandler{service: service}
}

func (h *RegisterIssuerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpjson.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	request, err := decodeJSONBody[models.RegisterIssuerRequest](r)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid JSON payload")
		return
	}

	response, err := h.service.RegisterIssuer(r.Context(), request)
	if err != nil {
		writeHandlerError(w, err)
		return
	}

	httpjson.WriteSuccess(w, http.StatusOK, response)
}
