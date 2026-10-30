package httpapi

import (
	"fmt"
	"net/http"
)

// cors wraps next with CORS headers for the allowed origins.
// It answers preflight requests itself with 204 and never passes them to next.
func cors(allowed []string, next http.Handler) http.Handler {
	allowedOrigins := make(map[string]bool)
	for _, origin := range allowed {
		allowedOrigins[origin] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		origin := r.Header.Get("Origin")
		fmt.Println(origin, allowedOrigins)
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Max-Age", "7200") // Chromium preflight cache maximum
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
