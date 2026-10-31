package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// cors returns middleware that sets CORS headers for the allowed origins.
// It answers preflight requests itself with 204 and aborts the chain.
func cors(allowed []string) gin.HandlerFunc {
	allowedOrigins := make(map[string]bool)
	for _, origin := range allowed {
		allowedOrigins[origin] = true
	}

	return func(c *gin.Context) {
		c.Writer.Header().Add("Vary", "Origin")
		origin := c.GetHeader("Origin")
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		if c.Request.Method == http.MethodOptions && c.GetHeader("Access-Control-Request-Method") != "" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
			c.Writer.Header().Set("Access-Control-Max-Age", "7200") // Chromium preflight cache maximum
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
