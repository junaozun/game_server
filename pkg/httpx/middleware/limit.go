package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit qps限制
func RateLimit(duration time.Duration, qps int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(duration), qps)
	return func(context *gin.Context) {
		if !limiter.Allow() {
			context.AbortWithStatus(http.StatusForbidden)
			return
		}
		context.Next()
	}
}
