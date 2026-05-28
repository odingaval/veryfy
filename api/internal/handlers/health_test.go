package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	handler := NewHealthHandler(time.Now())

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var body struct {
		Data struct {
			Status  string `json:"status"`
			Service string `json:"service"`
		} `json:"data"`
		Error any `json:"error"`
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Data.Status != "ok" {
		t.Fatalf("status payload = %q, want %q", body.Data.Status, "ok")
	}

	if body.Data.Service != "veryfy-api" {
		t.Fatalf("service payload = %q, want %q", body.Data.Service, "veryfy-api")
	}

	if body.Error != nil {
		t.Fatalf("error payload = %v, want nil", body.Error)
	}
}
