package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	keyAddr            = "ADDR"
	keyAllowedOrigins  = "ALLOWED_ORIGINS"
	keyPostsURL        = "POSTS_URL"
	keyRefreshInterval = "REFRESH_INTERVAL"

	defaultAddr            = "127.0.0.1:8080"
	defaultAllowedOrigins  = "http://localhost:4321"
	defaultPostsURL        = "https://ktblog.dev/posts.json"
	defaultRefreshInterval = 10 * time.Minute
)

type config struct {
	Addr            string
	AllowedOrigins  []string
	PostsURL        string
	RefreshInterval time.Duration
}

func loadConfig(getEnv func(string) string) (config, error) {
	addr := getEnv(keyAddr)
	if addr == "" {
		addr = defaultAddr
	}

	rawAllowedOrigins := getEnv(keyAllowedOrigins)
	if rawAllowedOrigins == "" {
		rawAllowedOrigins = defaultAllowedOrigins
	}
	allowedOrigins, err := parseOrigins(rawAllowedOrigins)
	if err != nil {
		return config{}, fmt.Errorf("parse %s: %w", keyAllowedOrigins, err)
	}
	if len(allowedOrigins) == 0 {
		return config{}, fmt.Errorf("%s must list at least one origin, got %q", keyAllowedOrigins, rawAllowedOrigins)
	}

	postsURL := getEnv(keyPostsURL)
	if postsURL == "" {
		postsURL = defaultPostsURL
	}

	var refreshInterval time.Duration
	rawInterval := getEnv(keyRefreshInterval)
	if rawInterval == "" {
		refreshInterval = defaultRefreshInterval
	} else {
		var err error
		refreshInterval, err = time.ParseDuration(rawInterval)
		if err != nil {
			return config{}, fmt.Errorf("parse %s: %w", keyRefreshInterval, err)
		}
		if refreshInterval <= 0 {
			return config{}, fmt.Errorf("%s must be positive, got %q", keyRefreshInterval, rawInterval)
		}
	}

	return config{Addr: addr, PostsURL: postsURL, RefreshInterval: refreshInterval, AllowedOrigins: allowedOrigins}, nil
}

func parseOrigins(raw string) ([]string, error) {
	allowedOrigins := make([]string, 0)
	for rawOrigin := range strings.SplitSeq(raw, ",") {
		origin := strings.TrimSpace(rawOrigin)
		if origin == "" {
			continue
		}
		if strings.HasSuffix(origin, "/") {
			return nil, fmt.Errorf("invalid origin %q: not expected with trailing slash", origin)
		}

		originURL, err := url.Parse(origin)
		if err != nil {
			return nil, fmt.Errorf("invalid origin %q: %w", origin, err)
		}
		if originURL.Scheme == "" {
			return nil, fmt.Errorf("invalid origin %q: missing scheme", origin)
		}
		if originURL.Host == "" {
			return nil, fmt.Errorf("invalid origin %q: missing host", origin)
		}
		if strings.ToLower(originURL.Host) != originURL.Host {
			return nil, fmt.Errorf("invalid origin %q: host should be lowercase", origin)
		}

		formURL := url.URL{Scheme: originURL.Scheme, Host: originURL.Host}
		if origin != formURL.String() {
			return nil, fmt.Errorf("invalid origin %q: want scheme://host[:port]", origin)
		}
		allowedOrigins = append(allowedOrigins, origin)
	}
	return allowedOrigins, nil
}
