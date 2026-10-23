package main

import (
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

			if tt.want != got {
				t.Errorf("loadConfig() = %+v, want = %+v", got, tt.want)
			}
		})
	}
}
