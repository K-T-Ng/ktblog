package httpapi_test

import (
	"encoding/json"
	"maps"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/K-T-Ng/ktblog/api/internal/httpapi"
	"github.com/K-T-Ng/ktblog/api/internal/stats"
)

type fakeAllowlist map[string]bool

func (f fakeAllowlist) Contains(slug string) bool {
	return f[slug]
}

func newTestHandler() http.Handler {
	return httpapi.NewHandler(httpapi.Deps{
		Stats:     stats.NewStore(),
		Allowlist: fakeAllowlist{"hello": true},
	})
}

func do(t *testing.T, h http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	h.ServeHTTP(rec, req)
	return rec
}

func assertHTTPStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if got := rec.Code; got != want {
		t.Fatalf("got HTTP status code %d, want %d", got, want)
	}
}

func assertJSONContentType(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	rawContentType := rec.Result().Header.Get("Content-Type")
	contentType, _, err := mime.ParseMediaType(rawContentType)
	if err != nil {
		t.Fatalf("parse Content-Type %q: %v", rawContentType, err)
	}
	if contentType != "application/json" {
		t.Fatalf("expected application/json content type, but got: %q", contentType)
	}
}

func assertCounts(t *testing.T, rec *httptest.ResponseRecorder, want map[string]int64) {
	t.Helper()
	assertHTTPStatus(t, rec, http.StatusOK)
	assertJSONContentType(t, rec)

	var body map[string]int64
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("expected response body is countsBody, but can't decode: %v", err)
	}
	if !maps.Equal(body, want) {
		t.Errorf("counts = %+v, want = %+v", body, want)
	}
}

func asInt(v any) int {
	if f, ok := v.(float64); ok {
		return int(f)
	}
	return 0
}

func assertAIPError(t *testing.T, rec *httptest.ResponseRecorder, wantCode int, wantStatus string) {
	t.Helper()
	assertHTTPStatus(t, rec, wantCode)
	assertJSONContentType(t, rec)

	var body map[string]map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("expected response body is errorBody, but can't decode: %v", err)
	}
	if asInt(body["error"]["code"]) != wantCode {
		t.Errorf("error code = %v, want %d", body["error"]["code"], wantCode)
	}
	if body["error"]["status"] != wantStatus {
		t.Errorf("error status = %q, want %q", body["error"]["status"], wantStatus)
	}
	if msg, _ := body["error"]["message"].(string); msg == "" {
		t.Errorf("expected AIP error message to be non-empty but got empty")
	}
}

func TestHealthz(t *testing.T) {
	// Arrange
	handler := newTestHandler()

	// Act
	rec := do(t, handler, http.MethodGet, "/healthz")

	// Assert
	assertHTTPStatus(t, rec, http.StatusOK)
	assertJSONContentType(t, rec)

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
	handler := newTestHandler()

	// Act
	rec := do(t, handler, http.MethodPost, "/healthz")

	// Assert
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestView(t *testing.T) {
	// Arrange
	handler := newTestHandler()

	// Act
	do(t, handler, http.MethodPost, "/v1/posts/hello:view")
	rec := do(t, handler, http.MethodPost, "/v1/posts/hello:view")

	// Assert
	assertCounts(t, rec, map[string]int64{"views": 2, "likes": 0})
}

func TestLikeUnlike(t *testing.T) {
	// Arrange
	handler := newTestHandler()

	// Act
	do(t, handler, http.MethodPost, "/v1/posts/hello:like")
	do(t, handler, http.MethodPost, "/v1/posts/hello:like")
	rec := do(t, handler, http.MethodPost, "/v1/posts/hello:unlike")

	// Assert
	assertCounts(t, rec, map[string]int64{"views": 0, "likes": 1})
}

func TestUnlikeAtZero(t *testing.T) {
	// Arrange
	handler := newTestHandler()

	// Act
	rec := do(t, handler, http.MethodPost, "/v1/posts/hello:unlike")

	// Assert
	assertCounts(t, rec, map[string]int64{"views": 0, "likes": 0})
}

func TestUnknownSlug(t *testing.T) {
	// Arrange
	handler := newTestHandler()

	// Act
	rec := do(t, handler, http.MethodPost, "/v1/posts/nope:view")

	// Assert
	assertAIPError(t, rec, http.StatusNotFound, "NOT_FOUND")
}

func TestUnknownMethod(t *testing.T) {
	// Arrange
	handler := newTestHandler()

	// Act
	rec := do(t, handler, http.MethodPost, "/v1/posts/hello:nope")

	// Assert
	assertAIPError(t, rec, http.StatusNotFound, "NOT_FOUND")
}

func TestNoMethod(t *testing.T) {
	// Arrange
	handler := newTestHandler()

	// Act
	rec := do(t, handler, http.MethodPost, "/v1/posts/hello")

	// Assert
	assertAIPError(t, rec, http.StatusNotFound, "NOT_FOUND")
}

func TestWrongHTTPMethod(t *testing.T) {
	// Arrange
	handler := newTestHandler()

	// Act
	rec := do(t, handler, http.MethodGet, "/v1/posts/hello:like")

	// Assert
	assertHTTPStatus(t, rec, http.StatusMethodNotAllowed)
}

func TestUnknownPath(t *testing.T) {
	// Arrange
	handler := newTestHandler()

	// Act
	rec := do(t, handler, http.MethodPost, "/nope")

	// Assert
	assertHTTPStatus(t, rec, http.StatusNotFound)
}

func TestBadRequestsDontCount(t *testing.T) {
	// Arrange
	handler := newTestHandler()

	// Act
	do(t, handler, http.MethodPost, "/v1/posts/hello:dance")
	rec := do(t, handler, http.MethodPost, "/v1/posts/hello:view")

	// Assert
	assertCounts(t, rec, map[string]int64{"views": 1, "likes": 0})
}
