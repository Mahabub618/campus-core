package router

import (
	"campus-core/internal/handler"
	"campus-core/internal/middleware"
	"campus-core/internal/repository"
	"campus-core/internal/service"

	"github.com/gin-gonic/gin"
)

func (r *Router) setupRoleRoutes(rg *gin.RouterGroup) {
	// Repositories
	userRepo := repository.NewUserRepository(r.db)
	teacherRepo := repository.NewTeacherRepository(r.db)
	studentRepo := repository.NewStudentRepository(r.db)
	parentRepo := repository.NewParentRepository(r.db)
	accountantRepo := repository.NewAccountantRepository(r.db)

	// Services
	teacherService := service.NewTeacherService(teacherRepo, userRepo, r.db, r.jwtManager)
	studentService := service.NewStudentService(studentRepo, userRepo, r.db, r.jwtManager)
	parentService := service.NewParentService(parentRepo, userRepo, r.db, r.jwtManager)
	accountantService := service.NewAccountantService(accountantRepo, userRepo, r.db, r.jwtManager)

	// Handlers
	teacherHandler := handler.NewTeacherHandler(teacherService)
	studentHandler := handler.NewStudentHandler(studentService)
	parentHandler := handler.NewParentHandler(parentService)
	accountantHandler := handler.NewAccountantHandler(accountantService)

	// Admin access required for creating roles (can be refined to RequirePermission)
	adminOnly := rg.Group("")
	adminOnly.Use(middleware.RequireAdmin())

	// Teachers
	teachers := adminOnly.Group("/teachers")
	{
		teachers.POST("", teacherHandler.Create)
		teachers.GET("", teacherHandler.GetAll)
		teachers.GET("/:id", teacherHandler.GetByID)
		teachers.PUT("/:id", teacherHandler.Update)
		teachers.GET("/:id/classes", teacherHandler.GetClasses)
		teachers.GET("/:id/subjects", teacherHandler.GetSubjects)
	}

	// Students
	students := adminOnly.Group("/students")
	{
		students.POST("", studentHandler.Create)
		students.GET("", studentHandler.GetAll)
		students.GET("/:id", studentHandler.GetByID)
		students.PUT("/:id", studentHandler.Update)
		students.GET("/:id/parents", studentHandler.GetParents)
		students.POST("/:id/parents", studentHandler.LinkParent)
		students.DELETE("/:id/parents/:parentId", studentHandler.UnlinkParent)
	}

	// Parents
	parents := adminOnly.Group("/parents")
	{
		parents.POST("", parentHandler.Create)
		parents.GET("", parentHandler.GetAll)
		parents.GET("/:id", parentHandler.GetByID)
		parents.PUT("/:id", parentHandler.Update)
		parents.GET("/:id/children", parentHandler.GetChildren)
	}

	// Accountants
	accountants := adminOnly.Group("/accountants")
	{
		accountants.POST("", accountantHandler.Create)
		accountants.GET("", accountantHandler.GetAll)
		accountants.GET("/:id", accountantHandler.GetByID)
		accountants.PUT("/:id", accountantHandler.Update)
	}
}
