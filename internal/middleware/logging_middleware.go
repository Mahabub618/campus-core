package middleware

import (
	"time"

	"campus-core/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestLogger returns a middleware that logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Build log fields
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.Int("status", statusCode),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.Int("body_size", c.Writer.Size()),
		}

		// Add user ID if available
		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}

		// Add institution ID if available
		if institutionID, exists := c.Get("institution_id"); exists {
			fields = append(fields, zap.Any("institution_id", institutionID))
		}

		// Log based on status code
		switch {
		case statusCode >= 500:
			logger.Error("Server error", fields...)
		case statusCode >= 400:
			logger.Warn("Client error", fields...)
		case statusCode >= 300:
			logger.Info("Redirect", fields...)
		default:
			logger.Info("Request completed", fields...)
		}
	}
}

// DebugLogger returns a more verbose logger for development
func DebugLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		start := time.Now()

		logger.Debug("Request started",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)

		c.Next()

		logger.Debug("Request completed",
			zap.String("request_id", requestID),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}
