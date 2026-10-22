package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/K-T-Ng/ktblog/api/internal/httpapi"
)

func TestHealthz(t *testing.T) {
	// Arrange
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	// Act
	httpapi.NewHandler().ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal health response: %v", err)
	}
	status, ok := resp["status"]
	if !ok {
		t.Errorf("expected 'status' field in health response")
	}
	if status != "ok" {
		t.Errorf("status = %s, want %s", status, "ok")
	}
}

func TestHealthzWrongMethod(t *testing.T) {
	// Arrange
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)

	// Act
	httpapi.NewHandler().ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
