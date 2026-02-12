package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	tokens float64
	last   time.Time
}

type RateLimiter struct {
	mu     sync.Mutex
	rate   float64
	burst  float64
	byIP   map[string]*bucket
	window time.Duration
}

func NewRateLimiter(ratePerSec, burst float64) *RateLimiter {
	return &RateLimiter{
		rate:   ratePerSec,
		burst: burst,
		byIP:  map[string]*bucket{},
		window: 10 * time.Minute,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := clientIP(c)
		now := time.Now()
		rl.mu.Lock()
		b, ok := rl.byIP[ip]
		if !ok {
			b = &bucket{tokens: rl.burst, last: now}
			rl.byIP[ip] = b
		}
		// refill
		elapsed := now.Sub(b.last).Seconds()
		b.tokens += elapsed * rl.rate
		if b.tokens > rl.burst {
			b.tokens = rl.burst
		}
		b.last = now

		allowed := b.tokens >= 1
		if allowed {
			b.tokens -= 1
		}
		rl.mu.Unlock()

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"ok": false, "error": "rate_limited"})
			return
		}
		c.Next()
	}
}

func clientIP(c *gin.Context) string {
	ip := c.ClientIP()
	// normalize
	if host, _, err := net.SplitHostPort(ip); err == nil {
		return host
	}
	return ip
}
