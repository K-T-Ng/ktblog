package httpapi

import (
	"fmt"
	"net/http"

	"github.com/K-T-Ng/ktblog/api/internal/stats"
	"github.com/gin-gonic/gin"
)

// NewHandler returns the API handler with all routes registered.
func NewHandler(d Deps) http.Handler {
	engine := gin.New()
	engine.HandleMethodNotAllowed = true
	engine.Use(gin.Recovery(), cors(d.AllowedOrigins))
	engine.GET("/healthz", healthHandler)
	engine.POST("/v1/posts/:name", postHandler(d))
	return engine
}

func healthHandler(c *gin.Context) {
	response := healthResponse{Status: "ok"}
	c.JSON(http.StatusOK, response)
}

// postHandler returns the handler for POST /v1/posts/{slug}:{method}.
// The custom method is checked before the slug, so a bad request never reaches Stats.
func postHandler(d Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse custom method, fail fast (404) if failed to parse
		customMethod := c.Param("name")
		slug, method, ok := parseCustomMethod(customMethod)
		if !ok {
			writeError(
				c, http.StatusNotFound, "NOT_FOUND",
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
			writeError(c, http.StatusNotFound, "NOT_FOUND", "unknown custom method")
			return
		}

		// Check slug exists
		if !d.Allowlist.Contains(slug) {
			writeError(c, http.StatusNotFound, "NOT_FOUND", "post not found")
			return
		}

		// Perform operations and response
		stat := op(slug)
		resp := countsResponse{Views: stat.Views, Likes: stat.Likes}
		c.JSON(http.StatusOK, resp)
	}
}
