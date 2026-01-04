package middleware

import (
	"strings"

	"campus-core/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware returns a middleware that validates JWT tokens
func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, 401, utils.ErrTokenMissing)
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.Error(c, 401, utils.ErrTokenInvalid)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			utils.Error(c, 401, err)
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_permissions", claims.Permissions)

		if claims.InstitutionID != "" {
			c.Set("institution_id", claims.InstitutionID)
		}

		c.Next()
	}
}

// OptionalAuthMiddleware returns a middleware that validates JWT tokens but doesn't require them
func OptionalAuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		claims, err := jwtManager.ValidateAccessToken(parts[1])
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)
			c.Set("user_permissions", claims.Permissions)
			if claims.InstitutionID != "" {
				c.Set("institution_id", claims.InstitutionID)
			}
		}

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	if id, ok := userID.(uuid.UUID); ok {
		return id, true
	}
	return uuid.Nil, false
}

// GetUserRole extracts user role from context
func GetUserRole(c *gin.Context) string {
	role, _ := c.Get("user_role")
	if r, ok := role.(string); ok {
		return r
	}
	return ""
}

// GetUserEmail extracts user email from context
func GetUserEmail(c *gin.Context) string {
	email, _ := c.Get("user_email")
	if e, ok := email.(string); ok {
		return e
	}
	return ""
}

// GetInstitutionID extracts institution ID from context
func GetInstitutionID(c *gin.Context) string {
	institutionID, _ := c.Get("institution_id")
	if id, ok := institutionID.(string); ok {
		return id
	}
	return ""
}

// GetUserPermissions extracts user permissions from context
func GetUserPermissions(c *gin.Context) []string {
	permissions, _ := c.Get("user_permissions")
	if p, ok := permissions.([]string); ok {
		return p
	}
	return []string{}
}
