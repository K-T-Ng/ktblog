package httpapi

import "github.com/K-T-Ng/ktblog/api/internal/stats"

// Stats counts views and likes per post slug.
// Each method returns the counts after the change.
type Stats interface {
	View(slug string) stats.Counts
	Like(slug string) stats.Counts
	Unlike(slug string) stats.Counts
}

// Allowlist reports which post slugs exist.
type Allowlist interface {
	Contains(slug string) bool
}

// Deps holds the dependencies of the handler returned by NewHandler.
type Deps struct {
	Stats          Stats
	Allowlist      Allowlist
	AllowedOrigins []string
}
