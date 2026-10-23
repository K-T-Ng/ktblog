package httpapi

import "testing"

func TestParseCustomMethod(t *testing.T) {
	// Arrange
	tests := []struct {
		name       string
		segment    string
		wantSlug   string
		wantMethod string
		wantOK     bool
	}{
		{
			name:       "valid semantic slug and method",
			segment:    "hello-world:like",
			wantSlug:   "hello-world",
			wantMethod: "like",
			wantOK:     true,
		},
		{
			name:       "valid random slug and semantic method",
			segment:    "a:view",
			wantSlug:   "a",
			wantMethod: "view",
			wantOK:     true,
		},
		{
			name:       "valid semantic slug and random method",
			segment:    "hello-world:b",
			wantSlug:   "hello-world",
			wantMethod: "b",
			wantOK:     true,
		},
		{
			name:       "multiple colons",
			segment:    "a:b:c",
			wantSlug:   "a",
			wantMethod: "b:c",
			wantOK:     true,
		},
		{
			name:       "no colon",
			segment:    "hello-world",
			wantSlug:   "",
			wantMethod: "",
			wantOK:     false,
		},
		{
			name:       "no method",
			segment:    "a:",
			wantSlug:   "",
			wantMethod: "",
			wantOK:     false,
		},
		{
			name:       "no slug",
			segment:    ":b",
			wantSlug:   "",
			wantMethod: "",
			wantOK:     false,
		},
		{
			name:       "empty input",
			segment:    "",
			wantSlug:   "",
			wantMethod: "",
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			slug, method, ok := parseCustomMethod(tt.segment)

			// Assert
			if slug != tt.wantSlug || method != tt.wantMethod || ok != tt.wantOK {
				t.Errorf(
					"parseCustomMethod() = (%q, %q, %t), want (%q, %q, %t)",
					slug, method, ok,
					tt.wantSlug, tt.wantMethod, tt.wantOK,
				)
			}
		})
	}
}
