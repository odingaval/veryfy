package httpjson

import (
	"encoding/json"
	"net/http"
)

// ErrorBody is the stable error shape returned by the API.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type envelope struct {
	Data  any        `json:"data"`
	Error *ErrorBody `json:"error"`
}

// WriteSuccess writes a successful JSON response using the API envelope.
func WriteSuccess(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, envelope{
		Data:  data,
		Error: nil,
	})
}

// WriteError writes an error JSON response using the API envelope.
func WriteError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, envelope{
		Data: nil,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
