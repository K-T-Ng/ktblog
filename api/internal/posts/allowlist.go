// Package posts tracks which post slugs exist on the site.
package posts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// Allowlist is the set of known post slugs, refreshed from posts.json.
// It is safe for concurrent use.
type Allowlist struct {
	url    string
	client *http.Client

	mu    sync.RWMutex
	slugs map[string]struct{}
}

type slugResponse struct {
	Slugs []string `json:"slugs"`
}

// NewAllowlist returns a new Allowlist.
func NewAllowlist(url string, client *http.Client) *Allowlist {
	return &Allowlist{
		url:    url,
		client: client,
		slugs:  make(map[string]struct{}),
	}
}

// Contains reports whether slug is a known post.
func (a *Allowlist) Contains(slug string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	_, ok := a.slugs[slug]
	return ok
}

// Refresh replaces the set with the slugs fetched from the URL. On any error, including an empty list,
// it keeps the current set.
func (a *Allowlist) Refresh(ctx context.Context) error {
	// Query slugs from web
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.url, nil)
	if err != nil {
		return fmt.Errorf("create request on GET %q: %w", a.url, err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("executing GET %q: %w", a.url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("close response body: %v", err)
		}
	}(resp.Body)

	// Parse slugs response from web
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code (%d) for GET %q", resp.StatusCode, a.url)
	}

	body := slugResponse{}
	decoder := json.NewDecoder(io.LimitReader(resp.Body, 1<<20))
	if err := decoder.Decode(&body); err != nil {
		return fmt.Errorf("decode response body for GET %q: %w", a.url, err)
	}

	if len(body.Slugs) == 0 {
		a.mu.RLock()
		slugLen := len(a.slugs)
		a.mu.RUnlock()
		return fmt.Errorf("empty slugs response from GET %q; keeping the previous %d", a.url, slugLen)
	}

	// Update available slugs
	newSlug := make(map[string]struct{})
	for _, slug := range body.Slugs {
		newSlug[slug] = struct{}{}
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.slugs = newSlug
	return nil
}

// Run refreshes the allowlist every interval until ctx is canceled.
// It does not refresh immediately, so call Refresh once before starting it.
// Refresh errors are logged, and the loop keeps running.
func (a *Allowlist) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := a.Refresh(ctx); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				log.Printf("refresh allowlist: %v", err)
			}
		}
	}
}
