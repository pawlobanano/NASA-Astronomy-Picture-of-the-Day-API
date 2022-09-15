package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// errorResponse returns a formatted error message.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// apiRateLimitExceededResponse returns a formatted error message.
func apiRateLimitExceededResponse() error {
	message := "API rate limit exceeded"
	return fmt.Errorf("%d %s", http.StatusTooManyRequests, message)
}
