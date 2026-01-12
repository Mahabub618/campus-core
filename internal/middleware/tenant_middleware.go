package middleware

import (
	"net/http"

	"campus-core/internal/database"
	"campus-core/internal/models"
	"campus-core/internal/utils"
	"campus-core/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TenantMiddleware handles multi-tenancy resolution
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check if institution_id is already in context (from AuthMiddleware)
		if authInstitutionID := GetInstitutionID(c); authInstitutionID != "" {
			// If header is also present, ensure it matches (security check)
			headerInstitutionID := c.GetHeader("X-Institution-ID")
			if headerInstitutionID != "" && headerInstitutionID != authInstitutionID {
				// Special case: Super Admin might be impersonating or accessing another tenant
				if GetUserRole(c) == models.RoleSuperAdmin {
					// Allow switching context for Super Admin
					logger.Info("Super Admin switching tenant context",
						zap.String("from", authInstitutionID),
						zap.String("to", headerInstitutionID))
					c.Set("institution_id", headerInstitutionID)
				} else {
					utils.Error(c, http.StatusForbidden, utils.ErrCrossTenantAccess)
					c.Abort()
					return
				}
			}
			c.Next()
			return
		}

		// 2. If not authenticated or no institution in token (e.g. Super Admin or Public Public), check Header
		institutionID := c.GetHeader("X-Institution-ID")
		if institutionID == "" {
			// For public endpoints that require tenant context
			// We don't abort here because some endpoints might be truly global (like /health or /login)
			// It's up to the specific handler/repository to check if institution_id is required
			c.Next()
			return
		}

		// 3. Validate Institution ID format
		id, err := uuid.Parse(institutionID)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
			c.Abort()
			return
		}

		// 4. Validate existence (Optional: Cache this check)
		// For now, we'll assume it exists to avoid DB hit on every request,
		// or we can do a quick check if we have a cache.
		// Since Redis is available, we could cache valid institution IDs.
		if database.RedisClient != nil {
			ctx := c.Request.Context()
			cacheKey := "institution:exists:" + institutionID
			exists, _ := database.Exists(ctx, cacheKey)
			if !exists {
				// Double check DB if not in cache (or if cache expired)
				var count int64
				if err := database.DB.Model(&models.Institution{}).Where("id = ? AND is_active = ?", id, true).Count(&count).Error; err != nil {
					logger.Error("Failed to check institution existence", zap.Error(err))
				} else if count == 0 {
					utils.Error(c, http.StatusNotFound, utils.ErrInstitutionNotFound)
					c.Abort()
					return
				} else {
					// Cache for 1 hour
					_ = database.SetWithExpiry(ctx, cacheKey, "1", 3600*1000000000) // 1 hour
				}
			}
		}

		c.Set("institution_id", institutionID)
		c.Next()
	}
}

// RequireTenant requires standard tenant context to be present
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetInstitutionID(c) == "" {
			utils.Error(c, http.StatusBadRequest, utils.ErrInstitutionIDRequired)
			c.Abort()
			return
		}
		c.Next()
	}
}
