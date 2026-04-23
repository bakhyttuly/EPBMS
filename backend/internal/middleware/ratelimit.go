package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// tokenBucket holds the state for a single IP's rate limit.
type tokenBucket struct {
	tokens   float64
	lastSeen time.Time
	mu       sync.Mutex
}

// RateLimiter returns a simple token-bucket rate limiter keyed by client IP.
// rps is the maximum allowed requests per second; burst is the maximum burst size.
func RateLimiter(rps float64, burst float64) gin.HandlerFunc {
	buckets := make(map[string]*tokenBucket)
	var globalMu sync.Mutex

	return func(c *gin.Context) {
		ip := c.ClientIP()

		globalMu.Lock()
		bucket, exists := buckets[ip]
		if !exists {
			bucket = &tokenBucket{tokens: burst, lastSeen: time.Now()}
			buckets[ip] = bucket
		}
		globalMu.Unlock()

		bucket.mu.Lock()
		defer bucket.mu.Unlock()

		now := time.Now()
		elapsed := now.Sub(bucket.lastSeen).Seconds()
		bucket.lastSeen = now

		// Refill tokens based on elapsed time.
		bucket.tokens += elapsed * rps
		if bucket.tokens > burst {
			bucket.tokens = burst
		}

		if bucket.tokens < 1 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "too many requests, please slow down",
			})
			c.Abort()
			return
		}

		bucket.tokens--
		c.Next()
	}
}
