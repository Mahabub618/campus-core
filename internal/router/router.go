package router

import (
	"campus-core/internal/config"
	"campus-core/internal/handler"
	"campus-core/internal/middleware"
	"campus-core/internal/repository"
	"campus-core/internal/service"
	"campus-core/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Router holds the Gin engine and dependencies
type Router struct {
	engine     *gin.Engine
	config     *config.Config
	db         *gorm.DB
	jwtManager *utils.JWTManager
}

// NewRouter creates a new router instance
func NewRouter(cfg *config.Config, db *gorm.DB) *Router {
	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Create Gin engine
	engine := gin.New()

	// Create JWT manager
	jwtManager := utils.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

	return &Router{
		engine:     engine,
		config:     cfg,
		db:         db,
		jwtManager: jwtManager,
	}
}

// Setup configures all routes and middleware
func (r *Router) Setup() *gin.Engine {
	// Apply global middleware
	r.engine.Use(middleware.Recovery())
	r.engine.Use(middleware.RequestLogger())
	r.engine.Use(middleware.CORS())

	// Apply rate limiting if Redis is available
	r.engine.Use(middleware.RateLimit(middleware.RateLimitConfig{
		Requests: r.config.RateLimit.Requests,
		Duration: r.config.RateLimit.Duration,
		KeyFunc:  func(c *gin.Context) string { return "ratelimit:" + c.ClientIP() },
	}))

	// Health check endpoint (no auth required)
	r.engine.GET("/api/v1/health", r.healthCheck)

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		// Setup auth routes
		r.setupAuthRoutes(v1)

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			// Tenant middleware to resolve institution context
			protected.Use(middleware.TenantMiddleware())

			r.setupInstitutionRoutes(protected)
			r.setupUserRoutes(protected)
			r.setupRoleRoutes(protected)

			// Academic management routes
			setupAcademicRoutes(protected, r.db)
		}
	}

	return r.engine
}

// setupAuthRoutes configures authentication routes
func (r *Router) setupAuthRoutes(rg *gin.RouterGroup) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(r.db)

	// Initialize services
	authService := service.NewAuthService(userRepo, r.jwtManager)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Auth routes group
	auth := rg.Group("/auth")
	{
		// Public routes (with stricter rate limiting)
		auth.POST("/login", middleware.AuthRateLimit(), authHandler.Login)
		auth.POST("/refresh-token", authHandler.RefreshToken)
		auth.POST("/forgot-password", middleware.AuthRateLimit(), authHandler.ForgotPassword)
		auth.POST("/reset-password", middleware.AuthRateLimit(), authHandler.ResetPassword)

		// Protected routes
		authProtected := auth.Group("")
		authProtected.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			authProtected.POST("/register", middleware.RequireAdmin(), authHandler.Register)
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.POST("/change-password", authHandler.ChangePassword)
			authProtected.GET("/me", authHandler.GetMe)
		}
	}
}

// healthCheck handles health check requests
func (r *Router) healthCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := r.db.DB()
	if err != nil {
		utils.InternalServerError(c, "Database connection error")
		return
	}

	if err := sqlDB.Ping(); err != nil {
		utils.InternalServerError(c, "Database ping failed")
		return
	}

	utils.OK(c, "Server is healthy", gin.H{
		"status":   "healthy",
		"database": "connected",
	})
}

// GetEngine returns the Gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// GetJWTManager returns the JWT manager
func (r *Router) GetJWTManager() *utils.JWTManager {
	return r.jwtManager
}
