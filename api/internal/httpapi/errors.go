package httpapi

import (
	"github.com/gin-gonic/gin"
)

// errorResponse is the AIP-193 error body.
type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// writeError writes an AIP-193 error response, such as
// {"error": {"code": 404, "status": "NOT_FOUND", "message": "post not found"}}.
func writeError(c *gin.Context, code int, status, message string) {
	errResp := errorResponse{
		Error: errorDetail{
			Code:    code,
			Status:  status,
			Message: message,
		},
	}
	c.JSON(code, errResp)
}
