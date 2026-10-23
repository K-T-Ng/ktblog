package main

import (
	"fmt"
	"time"
)

const (
	keyAddr            = "ADDR"
	keyPostsURL        = "POSTS_URL"
	keyRefreshInterval = "REFRESH_INTERVAL"

	defaultAddr            = "127.0.0.1:8080"
	defaultPostsURL        = "https://ktblog.dev/posts.json"
	defaultRefreshInterval = 10 * time.Minute
)

type config struct {
	Addr            string
	PostsURL        string
	RefreshInterval time.Duration
}

func loadConfig(getEnv func(string) string) (config, error) {
	addr := getEnv(keyAddr)
	if addr == "" {
		addr = defaultAddr
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

	return config{Addr: addr, PostsURL: postsURL, RefreshInterval: refreshInterval}, nil
}
