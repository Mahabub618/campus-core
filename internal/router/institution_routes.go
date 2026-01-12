package router

import (
	"campus-core/internal/handler"
	"campus-core/internal/middleware"
	"campus-core/internal/repository"
	"campus-core/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) setupInstitutionRoutes(rg *gin.RouterGroup) {
	repo := repository.NewInstitutionRepository(r.db)
	svc := service.NewInstitutionService(repo)
	handler := handler.NewInstitutionHandler(svc)

	institutions := rg.Group("/institutions")
	// Only Super Admin can manage institutions
	institutions.Use(middleware.RequireSuperAdmin())
	{
		institutions.POST("", handler.Create)
		institutions.GET("", handler.GetAll)
		institutions.GET("/:id", handler.GetByID)
		institutions.PUT("/:id", handler.Update)
		institutions.DELETE("/:id", handler.Delete)
		institutions.GET("/:id/stats", handler.GetStats)
	}
}
