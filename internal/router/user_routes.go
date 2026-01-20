package router

import (
	"campus-core/internal/handler"
	"campus-core/internal/middleware"
	"campus-core/internal/repository"
	"campus-core/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) setupUserRoutes(rg *gin.RouterGroup) {
	// Repos
	userRepo := repository.NewUserRepository(r.db)
	instRepo := repository.NewInstitutionRepository(r.db)

	// Services
	// Note: We need existing AuthService instance, or create new one?
	// Router has jwtManager, but AuthService needs Repo + JWT.
	// Ideally we accept AuthService in router setup or create it.
	// In `router.go`, we created `authService` inside `setupAuthRoutes` locally.
	// We should probably promote `authService` to struct level or recreate (stateless except for repo).
	authService := service.NewAuthService(userRepo, r.jwtManager)
	userService := service.NewUserService(userRepo, instRepo, authService)
	userHandler := handler.NewUserHandler(userService)

	users := rg.Group("/users")
	users.Use(middleware.RequireAdmin()) // Only Admins can manage users
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.GetAllUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.PATCH("/:id/status", userHandler.ToggleStatus)
	}

	profile := rg.Group("/profile")
	// Any authenticated user can access profile
	{
		profile.GET("", userHandler.GetProfile)
		profile.PUT("", userHandler.UpdateProfile)
		profile.PUT("/avatar", userHandler.UpdateAvatar)
		profile.PUT("/password", userHandler.UpdatePassword)
	}
}
