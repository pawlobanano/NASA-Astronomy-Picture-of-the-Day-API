package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// apiRateLimiter limits the request rate on API.
func apiRateLimiter(apiRate int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(time.Second), apiRate)

	return func(ctx *gin.Context) {
		if !limiter.Allow() {
			ctx.JSON(http.StatusTooManyRequests, errorResponse(apiRateLimitExceededResponse()))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
