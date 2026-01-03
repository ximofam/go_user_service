package middleware

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type tokenBucket struct {
	capacity   float64
	tokens     float64
	rate       float64 // (token/s) thêm token sau 1 s
	lastRefill time.Time
	mu         *sync.Mutex
}

func (t *tokenBucket) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	elapse := now.Sub(t.lastRefill).Seconds()

	newTokens := elapse * t.rate
	t.tokens = math.Min(t.capacity, t.tokens+newTokens)
	t.lastRefill = now

	if t.tokens >= 1 {
		t.tokens--
		return true
	}

	return false
}

type ipRateLimiter struct {
	ips      map[string]*tokenBucket
	mu       *sync.Mutex
	capacity float64
	rate     float64
	ipTTL    time.Duration
}

func NewIPRateLimiter(capacity, rate float64, ipTTL time.Duration) *ipRateLimiter {
	ipRateLimiter := ipRateLimiter{
		ips:      map[string]*tokenBucket{},
		mu:       &sync.Mutex{},
		capacity: capacity,
		rate:     rate,
		ipTTL:    ipTTL,
	}

	ipRateLimiter.startCleanUp()

	return &ipRateLimiter
}

func (rl *ipRateLimiter) GetBucket(ip string) *tokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, ok := rl.ips[ip]
	if !ok {
		bucket = &tokenBucket{
			capacity:   rl.capacity,
			tokens:     rl.capacity,
			rate:       rl.rate,
			lastRefill: time.Now(),
			mu:         rl.mu,
		}
		rl.ips[ip] = bucket
	}

	return bucket
}

func (rl *ipRateLimiter) startCleanUp() {
	ticker := time.NewTicker(rl.ipTTL)

	go func() {
		for now := range ticker.C {
			rl.cleanUp(now)
		}
	}()
}

func (rl *ipRateLimiter) cleanUp(now time.Time) {
	for ip, bucket := range rl.ips {
		rl.mu.Lock()
		if now.Sub(bucket.lastRefill) > rl.ipTTL {
			delete(rl.ips, ip)
			log.Printf("RateLimter ip: %s", ip)
		}
		rl.mu.Unlock()
	}
}

func RateLimiterMiddleware(limiter *ipRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		client := limiter.GetBucket(ip)

		if !client.Allow() {
			c.AbortWithStatusJSON(429, gin.H{
				"error": "Too many requests",
			})
			return
		}
		c.Next()
	}
}
