package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"campus-core/internal/database"
	"campus-core/internal/utils"
	"campus-core/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int                       // Maximum number of requests
	Duration time.Duration             // Time window
	KeyFunc  func(*gin.Context) string // Function to generate the rate limit key
}

// DefaultRateLimitConfig returns default rate limit config
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Requests: 100,
		Duration: 1 * time.Minute,
		KeyFunc:  defaultKeyFunc,
	}
}

// defaultKeyFunc uses client IP as the rate limit key
func defaultKeyFunc(c *gin.Context) string {
	return "ratelimit:" + c.ClientIP()
}

// userKeyFunc uses user ID as the rate limit key (for authenticated requests)
func UserKeyFunc(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("ratelimit:user:%v", userID)
	}
	return "ratelimit:" + c.ClientIP()
}

// RateLimit returns a rate limiting middleware
func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if database.RedisClient == nil {
			// Skip rate limiting if Redis is not available
			logger.Warn("Rate limiting skipped: Redis not connected")
			c.Next()
			return
		}

		ctx := context.Background()
		key := config.KeyFunc(c)

		// Get current count
		count, err := database.RedisClient.Get(ctx, key).Int64()
		if err != nil && err.Error() != "redis: nil" {
			logger.Error("Rate limit check failed", zap.Error(err))
			c.Next()
			return
		}

		// Check if limit exceeded
		if count >= int64(config.Requests) {
			// Get TTL for Retry-After header
			ttl, _ := database.RedisClient.TTL(ctx, key).Result()

			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.Requests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, utils.ErrorResponse{
				Success: false,
				Error:   "Rate limit exceeded. Please try again later.",
				Code:    "SYS_005",
			})
			return
		}

		// Increment counter
		pipe := database.RedisClient.Pipeline()
		pipe.Incr(ctx, key)

		// Set expiry only if key doesn't exist
		if count == 0 {
			pipe.Expire(ctx, key, config.Duration)
		}

		_, err = pipe.Exec(ctx)
		if err != nil {
			logger.Error("Rate limit increment failed", zap.Error(err))
		}

		// Set rate limit headers
		remaining := config.Requests - int(count) - 1
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.Requests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		c.Next()
	}
}

// StrictRateLimit returns a stricter rate limit for sensitive endpoints
func StrictRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Requests: 10,
		Duration: 1 * time.Minute,
		KeyFunc:  defaultKeyFunc,
	})
}

// AuthRateLimit returns rate limiting for auth endpoints (login, password reset)
func AuthRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Requests: 5,
		Duration: 1 * time.Minute,
		KeyFunc:  defaultKeyFunc,
	})
}
