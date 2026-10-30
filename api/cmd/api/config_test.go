package main

import (
	"reflect"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		env     map[string]string
		want    config
		wantErr bool
	}{
		{
			name: "defaults",
			env:  map[string]string{},
			want: config{
				Addr:            "127.0.0.1:8080",
				AllowedOrigins:  []string{"http://localhost:4321"},
				PostsURL:        "https://ktblog.dev/posts.json",
				RefreshInterval: 10 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "addr override",
			env:  map[string]string{"ADDR": "127.0.0.1:9090"},
			want: config{
				Addr:            "127.0.0.1:9090",
				AllowedOrigins:  []string{"http://localhost:4321"},
				PostsURL:        "https://ktblog.dev/posts.json",
				RefreshInterval: 10 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "posts url override",
			env:  map[string]string{"POSTS_URL": "http://localhost:4321/posts.json"},
			want: config{
				Addr:            "127.0.0.1:8080",
				AllowedOrigins:  []string{"http://localhost:4321"},
				PostsURL:        "http://localhost:4321/posts.json",
				RefreshInterval: 10 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "interval override",
			env:  map[string]string{"REFRESH_INTERVAL": "30s"},
			want: config{
				Addr:            "127.0.0.1:8080",
				AllowedOrigins:  []string{"http://localhost:4321"},
				PostsURL:        "https://ktblog.dev/posts.json",
				RefreshInterval: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name:    "interval garbage",
			env:     map[string]string{"REFRESH_INTERVAL": "often"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "interval no unit",
			env:     map[string]string{"REFRESH_INTERVAL": "30"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "interval zero",
			env:     map[string]string{"REFRESH_INTERVAL": "0"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "interval negative",
			env:     map[string]string{"REFRESH_INTERVAL": "-1m"},
			want:    config{},
			wantErr: true,
		},
		{
			name: "two allowed origins separated with spaces",
			env:  map[string]string{"ALLOWED_ORIGINS": "https://ktblog.dev, http://localhost:4321"},
			want: config{
				Addr:            "127.0.0.1:8080",
				AllowedOrigins:  []string{"https://ktblog.dev", "http://localhost:4321"},
				PostsURL:        "https://ktblog.dev/posts.json",
				RefreshInterval: 10 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "allowed origins with trailing comma",
			env:  map[string]string{"ALLOWED_ORIGINS": "https://ktblog.dev,"},
			want: config{
				Addr:            "127.0.0.1:8080",
				AllowedOrigins:  []string{"https://ktblog.dev"},
				PostsURL:        "https://ktblog.dev/posts.json",
				RefreshInterval: 10 * time.Minute,
			},
			wantErr: false,
		},
		{
			name:    "allowed origins with trailing slash",
			env:     map[string]string{"ALLOWED_ORIGINS": "https://ktblog.dev/"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "allowed origins with scheme missing",
			env:     map[string]string{"ALLOWED_ORIGINS": "ktblog.dev"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "allowed origins with path",
			env:     map[string]string{"ALLOWED_ORIGINS": "https://ktblog.dev/blog"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "allowed origins with empty raw query",
			env:     map[string]string{"ALLOWED_ORIGINS": "https://ktblog.dev?"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "allowed origins with empty host",
			env:     map[string]string{"ALLOWED_ORIGINS": "https:"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "allowed origins with user set",
			env:     map[string]string{"ALLOWED_ORIGINS": "https://user@ktblog.dev"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "allowed origins with scheme uppercase",
			env:     map[string]string{"ALLOWED_ORIGINS": "HTTPS://ktblog.dev"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "allowed origins with host uppercase",
			env:     map[string]string{"ALLOWED_ORIGINS": "https://KTBLOG.DEV"},
			want:    config{},
			wantErr: true,
		},
		{
			name:    "allowed origins set but empty",
			env:     map[string]string{"ALLOWED_ORIGINS": ","},
			want:    config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got, err := loadConfig(func(k string) string { return tt.env[k] })

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			} else if !tt.wantErr && err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("loadConfig() = %+v, want = %+v", got, tt.want)
			}
		})
	}
}
