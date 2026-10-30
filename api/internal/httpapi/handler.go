package httpapi

import (
	"fmt"
	"net/http"

	"github.com/K-T-Ng/ktblog/api/internal/stats"
)

// NewHandler returns the API handler with all routes registered.
func NewHandler(d Deps) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("POST /v1/posts/{name}", postHandler(d))
	return cors(d.AllowedOrigins, mux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := healthResponse{Status: "ok"}
	writeJSON(w, http.StatusOK, response)
}

// postHandler returns the handler for POST /v1/posts/{slug}:{method}.
// The custom method is checked before the slug, so a bad request never reaches Stats.
func postHandler(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse custom method, fail fast (404) if failed to parse
		customMethod := r.PathValue("name")
		slug, method, ok := parseCustomMethod(customMethod)
		if !ok {
			writeError(
				w, http.StatusNotFound, "NOT_FOUND",
				fmt.Sprintf("cannot parse custom method %q", customMethod),
			)
			return
		}

		// Pick operation by method
		var op func(string) stats.Counts
		switch method {
		case "view":
			op = d.Stats.View
		case "like":
			op = d.Stats.Like
		case "unlike":
			op = d.Stats.Unlike
		default:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "unknown custom method")
			return
		}

		// Check slug exists
		if !d.Allowlist.Contains(slug) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "post not found")
			return
		}

		// Perform operations and response
		stat := op(slug)
		resp := countsResponse{Views: stat.Views, Likes: stat.Likes}
		writeJSON(w, http.StatusOK, resp)
	}
}
