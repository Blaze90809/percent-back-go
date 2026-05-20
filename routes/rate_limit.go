package routes

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// GeneralLimiter allows 120 requests per minute with a burst of 60
	GeneralLimiter = NewIPRateLimiter(2.0, 60.0)
	// SensitiveLimiter allows 5 requests per minute with a burst of 5
	SensitiveLimiter = NewIPRateLimiter(0.0833, 5.0)
)

type limiter struct {
	tokens     float64
	lastRefill time.Time
}

type IPRateLimiter struct {
	mu    sync.Mutex
	ips   map[string]*limiter
	rate  float64 // tokens refilled per second
	burst float64 // maximum capacity of bucket
}

// NewIPRateLimiter creates a new thread-safe rate limiter.
func NewIPRateLimiter(r float64, b float64) *IPRateLimiter {
	limiterInstance := &IPRateLimiter{
		ips:   make(map[string]*limiter),
		rate:  r,
		burst: b,
	}
	// Periodically clean up stale IPs in the background to prevent memory growth
	go limiterInstance.cleanUpLoop(10 * time.Minute)
	return limiterInstance
}

func (l *IPRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	lim, exists := l.ips[ip]
	if !exists {
		// New IP gets burst-1 tokens since they just consumed 1 token
		l.ips[ip] = &limiter{
			tokens:     l.burst - 1.0,
			lastRefill: now,
		}
		return true
	}

	// Refill tokens based on time elapsed
	elapsed := now.Sub(lim.lastRefill).Seconds()
	lim.tokens += elapsed * l.rate
	if lim.tokens > l.burst {
		lim.tokens = l.burst
	}
	lim.lastRefill = now

	// Consume 1 token if available
	if lim.tokens >= 1.0 {
		lim.tokens -= 1.0
		return true
	}

	return false
}

func (l *IPRateLimiter) cleanUpLoop(interval time.Duration) {
	for {
		time.Sleep(interval)
		l.mu.Lock()
		now := time.Now()
		for ip, lim := range l.ips {
			// If the IP hasn't been active for the cleanup interval, remove it
			if now.Sub(lim.lastRefill) > interval {
				delete(l.ips, ip)
			}
		}
		l.mu.Unlock()
	}
}

// RateLimitMiddleware returns a Gin middleware that rate-limits based on client IP.
func RateLimitMiddleware(limiterInstance *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiterInstance.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down and try again.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
