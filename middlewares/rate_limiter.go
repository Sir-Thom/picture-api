package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu    sync.Mutex
	ips   map[string]*ipLimiter
	rate  rate.Limit
	burst int
	ttl   time.Duration
}

func NewRateLimiter(r rate.Limit, b int, cleanupInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		ips:   make(map[string]*ipLimiter),
		rate:  r,
		burst: b,
		ttl:   cleanupInterval,
	}

	// Periodically clean up idle IP limiters to prevent memory leaks
	go func() {
		for {
			time.Sleep(cleanupInterval)
			rl.mu.Lock()
			for ip, limiter := range rl.ips {
				if time.Since(limiter.lastSeen) > cleanupInterval {
					delete(rl.ips, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		lim, exists := rl.ips[ip]
		if !exists {
			lim = &ipLimiter{
				limiter: rate.NewLimiter(rl.rate, rl.burst),
			}
			rl.ips[ip] = lim
		}
		lim.lastSeen = time.Now()
		rl.mu.Unlock()

		if !lim.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
