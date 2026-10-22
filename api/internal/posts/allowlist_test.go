package posts_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/K-T-Ng/ktblog/api/internal/posts"
)

// fakePosts serves whatever status and body are currently set.
type fakePosts struct {
	mu     sync.Mutex
	status int
	body   string
}

func (f *fakePosts) set(status int, body string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.status = status
	f.body = body
}

func (f *fakePosts) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.WriteHeader(f.status)
	if _, err := w.Write([]byte(f.body)); err != nil {
		log.Printf("failed to write fakePosts response: %v", err)
	}
}

func newFakeServer(t *testing.T, status int, body string) (*fakePosts, *httptest.Server) {
	f := &fakePosts{status: status, body: body}
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	return f, srv
}

func TestRefreshLoadsSlugs(t *testing.T) {
	// Arrange
	_, srv := newFakeServer(t, http.StatusOK, `{"slugs": ["a", "b"]}`)

	al := posts.NewAllowlist(srv.URL, &http.Client{})

	// Act
	if err := al.Refresh(context.Background()); err != nil {
		t.Fatalf("failed to refresh allowlist: %v", err)
	}

	// Assert
	if !al.Contains("a") || !al.Contains("b") {
		t.Errorf("expected slug a and b in allowlist")
	}
	if al.Contains("c") {
		t.Errorf("unexpected slug c in allowlist")
	}
}

func TestContainsBeforeRefresh(t *testing.T) {
	// Arrange
	al := posts.NewAllowlist("", &http.Client{})

	// Act
	contains := al.Contains("a")

	// Assert
	if contains {
		t.Errorf("unexpected a in allowlist")
	}
}

func TestRefreshReplaceSet(t *testing.T) {
	// Arrange
	f, srv := newFakeServer(t, http.StatusOK, `{"slugs": ["a"]}`)

	al := posts.NewAllowlist(srv.URL, &http.Client{})

	// Act
	if err := al.Refresh(context.Background()); err != nil {
		t.Fatalf("failed to refresh allowlist: %v", err)
	}

	f.set(http.StatusOK, `{"slugs": ["b"]}`)
	if err := al.Refresh(context.Background()); err != nil {
		t.Fatalf("failed to refresh allowlist: %v", err)
	}

	// Assert
	if al.Contains("a") {
		t.Errorf("unexpected a in allowlist")
	}
	if !al.Contains("b") {
		t.Errorf("expected b in allowlist")
	}
}

func TestRefreshNon200KeepsOld(t *testing.T) {
	// Arrange
	f, srv := newFakeServer(t, http.StatusOK, `{"slugs": ["a"]}`)

	al := posts.NewAllowlist(srv.URL, &http.Client{})

	// Act
	if err := al.Refresh(context.Background()); err != nil {
		t.Fatalf("failed to refresh allowlist: %v", err)
	}

	f.set(http.StatusServiceUnavailable, `<html></html>`)
	if err := al.Refresh(context.Background()); err == nil {
		t.Fatalf("expected error non-nil when non-200 response")
	}

	// Assert
	if !al.Contains("a") {
		t.Errorf("expected a in allowlist")
	}
}

func TestRefreshBadJSONKeepsOld(t *testing.T) {
	// Arrange
	f, srv := newFakeServer(t, http.StatusOK, `{"slugs": ["a"]}`)

	al := posts.NewAllowlist(srv.URL, &http.Client{})

	// Act
	if err := al.Refresh(context.Background()); err != nil {
		t.Fatalf("failed to refresh allowlist: %v", err)
	}

	f.set(http.StatusOK, `{"not" a json}`)
	if err := al.Refresh(context.Background()); err == nil {
		t.Fatalf("expected error non-nil when non-JSON response")
	}

	// Assert
	if !al.Contains("a") {
		t.Errorf("expected a in allowlist")
	}
}

func TestRefreshEmptyKeepsOld(t *testing.T) {
	// Arrange
	f, srv := newFakeServer(t, http.StatusOK, `{"slugs": ["a"]}`)

	al := posts.NewAllowlist(srv.URL, &http.Client{})

	// Act
	if err := al.Refresh(context.Background()); err != nil {
		t.Fatalf("failed to refresh allowlist: %v", err)
	}

	f.set(http.StatusOK, `{"slugs": []}`)
	if err := al.Refresh(context.Background()); err == nil {
		t.Fatalf("expected empty list is error when refreshing, but got nil")
	}

	// Assert
	if !al.Contains("a") {
		t.Errorf("expected a in allowlist")
	}
}

func TestRefreshUnreachable(t *testing.T) {
	// Arrange
	_, srv := newFakeServer(t, http.StatusOK, "")
	srv.Close()

	al := posts.NewAllowlist(srv.URL, &http.Client{})

	// Act & Assert
	if err := al.Refresh(context.Background()); err == nil {
		t.Errorf("expected error non-nil when unreachable")
	}
}
