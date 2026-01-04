package middleware

import (
	"net/http"
	"runtime/debug"

	"campus-core/internal/utils"
	"campus-core/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery returns a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := string(debug.Stack())

				// Log the panic
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("stack", stack),
				)

				// Abort and return error response
				c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse{
					Success: false,
					Error:   "Internal server error",
					Code:    "SYS_001",
				})
			}
		}()

		c.Next()
	}
}

// RecoveryWithCallback returns a recovery middleware with a custom callback
func RecoveryWithCallback(callback func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())

				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("stack", stack),
				)

				// Call the callback
				if callback != nil {
					callback(c, err)
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrorResponse{
					Success: false,
					Error:   "Internal server error",
					Code:    "SYS_001",
				})
			}
		}()

		c.Next()
	}
}
