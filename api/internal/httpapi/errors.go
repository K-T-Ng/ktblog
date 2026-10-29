package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// errorResponse is the AIP-193 error body.
type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// writeJSON writes v as a JSON response with the given status code.
// Encoding errors are only logged: the status line has already been sent.
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encoding JSON response: %v", err)
	}
}

// writeError writes an AIP-193 error response, such as
// {"error": {"code": 404, "status": "NOT_FOUND", "message": "post not found"}}.
func writeError(w http.ResponseWriter, code int, status, message string) {
	errResp := errorResponse{
		Error: errorDetail{
			Code:    code,
			Status:  status,
			Message: message,
		},
	}
	writeJSON(w, code, errResp)
}
