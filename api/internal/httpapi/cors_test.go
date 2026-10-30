package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var corsForTest = []string{"https://ktblog.dev"}

func assertACAOEquals(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()
	actual := rec.Result().Header.Get("Access-Control-Allow-Origin")
	if actual != want {
		t.Errorf("got Access-Control-Allow-Origin %q, want %q", actual, want)
	}
}

func assertAllowMethodContains(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()
	actual := rec.Result().Header.Get("Access-Control-Allow-Methods")
	for method := range strings.SplitSeq(actual, ",") {
		trimmed := strings.TrimSpace(method)
		if strings.EqualFold(trimmed, want) {
			return
		}
	}
	t.Errorf("expected Access-Control-Allow-Methods contains %q", want)
}

func assertMaxAgeNotEmpty(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	maxAge := rec.Result().Header.Get("Access-Control-Max-Age")
	if maxAge == "" {
		t.Errorf("expected Access-Control-Max-Age header to be non-empty")
		return
	}
	if _, err := strconv.Atoi(maxAge); err != nil {
		t.Errorf("expected Access-Control-Max-Age %q to be an integer: %v", maxAge, err)
	}
}

func assertVaryContains(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()
	varySlice := rec.Result().Header.Values("Vary")
	for _, vary := range varySlice {
		for piece := range strings.SplitSeq(vary, ",") {
			trimmed := strings.TrimSpace(piece)
			if strings.EqualFold(trimmed, want) {
				return
			}
		}
	}
	t.Errorf("expected Vary contains %q", want)
}

func TestCORSAllowedOrigin(t *testing.T) {
	// Arrange
	handler := newTestHandler(corsForTest)

	// Act
	rec := do(t, handler, http.MethodPost, "/v1/posts/hello:view", map[string]string{"Origin": "https://ktblog.dev"})

	// Assert
	assertHTTPStatus(t, rec, http.StatusOK)
	assertACAOEquals(t, rec, "https://ktblog.dev")
	assertVaryContains(t, rec, "Origin")
}

func TestCORSOtherOrigin(t *testing.T) {
	// Arrange
	handler := newTestHandler(corsForTest)

	// Act
	rec := do(t, handler, http.MethodPost, "/v1/posts/hello:view", map[string]string{"Origin": "https://evil.example"})

	// Assert
	assertHTTPStatus(t, rec, http.StatusOK)
	assertACAOEquals(t, rec, "")
	assertVaryContains(t, rec, "Origin")
}

func TestCORSNoOrigin(t *testing.T) {
	// Arrange
	handler := newTestHandler(corsForTest)

	// Act
	rec := do(t, handler, http.MethodPost, "/v1/posts/hello:view", nil)

	// Assert
	assertHTTPStatus(t, rec, http.StatusOK)
	assertACAOEquals(t, rec, "")
}

func TestPreflightAllowed(t *testing.T) {
	// Arrange
	handler := newTestHandler(corsForTest)

	// Act
	rec := do(
		t, handler, http.MethodOptions, "/v1/posts/hello:like",
		map[string]string{"Origin": "https://ktblog.dev", "Access-Control-Request-Method": "POST"},
	)

	// Assert
	assertHTTPStatus(t, rec, http.StatusNoContent)
	assertACAOEquals(t, rec, "https://ktblog.dev")
	assertAllowMethodContains(t, rec, "POST")
	assertMaxAgeNotEmpty(t, rec)
	assertVaryContains(t, rec, "Origin")
}

func TestPreflightOtherOrigin(t *testing.T) {
	// Arrange
	handler := newTestHandler(corsForTest)

	// Act
	rec := do(
		t, handler, http.MethodOptions, "/v1/posts/hello:like",
		map[string]string{"Origin": "https://evil.example", "Access-Control-Request-Method": "POST"},
	)

	// Assert
	assertHTTPStatus(t, rec, http.StatusNoContent)
	assertACAOEquals(t, rec, "")
}

func TestPreflightDoesNotCount(t *testing.T) {
	// Arrange
	handler := newTestHandler(corsForTest)

	// Act
	do(
		t, handler, http.MethodOptions, "/v1/posts/hello:view",
		map[string]string{"Origin": "https://ktblog.dev", "Access-Control-Request-Method": "POST"},
	)
	rec := do(
		t, handler, http.MethodPost, "/v1/posts/hello:view",
		map[string]string{"Origin": "https://ktblog.dev"},
	)

	// Assert
	assertCounts(t, rec, map[string]int64{"views": 1, "likes": 0})
}
