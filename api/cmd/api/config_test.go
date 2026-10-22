package main

import "testing"

func TestLoadConfig(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		env     map[string]string
		want    config
		wantErr bool
	}{
		{
			name:    "defaults",
			env:     map[string]string{},
			want:    config{Addr: "127.0.0.1:8080"},
			wantErr: false,
		},
		{
			name:    "addr override",
			env:     map[string]string{"ADDR": "127.0.0.1:9090"},
			want:    config{Addr: "127.0.0.1:9090"},
			wantErr: false,
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
