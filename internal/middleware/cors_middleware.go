package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Institution-ID", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// CORS returns a CORS middleware handler
func CORS() gin.HandlerFunc {
	return CORSWithConfig(DefaultCORSConfig())
}

// CORSWithConfig returns a CORS middleware handler with custom config
func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowOrigin := ""
		for _, o := range config.AllowOrigins {
			if o == "*" || o == origin {
				allowOrigin = origin
				if o == "*" {
					allowOrigin = "*"
				}
				break
			}
		}

		if allowOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowOrigin)
			c.Header("Access-Control-Allow-Methods", joinStrings(config.AllowMethods))
			c.Header("Access-Control-Allow-Headers", joinStrings(config.AllowHeaders))
			c.Header("Access-Control-Expose-Headers", joinStrings(config.ExposeHeaders))

			if config.AllowCredentials && allowOrigin != "*" {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.Header("Access-Control-Max-Age", intToString(config.MaxAge))
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}

	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
