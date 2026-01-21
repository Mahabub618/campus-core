package router

import (
	"campus-core/internal/handler"
	"campus-core/internal/middleware"
	"campus-core/internal/repository"
	"campus-core/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupAcademicRoutes configures all academic management routes
func setupAcademicRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories
	academicYearRepo := repository.NewAcademicYearRepository(db)
	classRepo := repository.NewClassRepository(db)
	sectionRepo := repository.NewSectionRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)
	timetableRepo := repository.NewTimetableRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)

	// Initialize services
	academicYearService := service.NewAcademicYearService(academicYearRepo)
	classService := service.NewClassService(classRepo, sectionRepo, teacherRepo)
	subjectService := service.NewSubjectService(subjectRepo, classRepo, teacherRepo)
	departmentService := service.NewDepartmentService(departmentRepo, teacherRepo)
	timetableService := service.NewTimetableService(
		timetableRepo, classRepo, sectionRepo, subjectRepo, teacherRepo, academicYearRepo,
	)

	// Initialize handlers
	academicYearHandler := handler.NewAcademicYearHandler(academicYearService)
	classHandler := handler.NewClassHandler(classService)
	subjectHandler := handler.NewSubjectHandler(subjectService)
	departmentHandler := handler.NewDepartmentHandler(departmentService)
	timetableHandler := handler.NewTimetableHandler(timetableService)

	// Academic Years routes
	academicYears := rg.Group("/academic-years")
	{
		academicYears.GET("", academicYearHandler.GetAll)
		academicYears.GET("/current", academicYearHandler.GetCurrent)
		academicYears.GET("/:id", academicYearHandler.GetByID)

		// Admin only routes
		academicYears.POST("", middleware.RequireAdmin(), academicYearHandler.Create)
		academicYears.PUT("/:id", middleware.RequireAdmin(), academicYearHandler.Update)
		academicYears.PATCH("/:id/activate", middleware.RequireAdmin(), academicYearHandler.Activate)
		academicYears.DELETE("/:id", middleware.RequireAdmin(), academicYearHandler.Delete)
	}

	// Classes routes
	classes := rg.Group("/classes")
	{
		classes.GET("", classHandler.GetAll)
		classes.GET("/:id", classHandler.GetByID)
		classes.GET("/:id/students", classHandler.GetStudents)
		classes.GET("/:id/teachers", classHandler.GetTeachers)

		// Admin only routes
		classes.POST("", middleware.RequireAdmin(), classHandler.Create)
		classes.PUT("/:id", middleware.RequireAdmin(), classHandler.Update)
		classes.DELETE("/:id", middleware.RequireAdmin(), classHandler.Delete)
	}

	// Sections routes (nested under classes)
	sections := rg.Group("/classes/:id/sections")
	{
		sections.GET("", classHandler.GetSections)
		sections.POST("", middleware.RequireAdmin(), classHandler.CreateSection)
	}

	// Standalone section routes
	sectionRoutes := rg.Group("/sections")
	{
		sectionRoutes.GET("/:id/students", classHandler.GetSectionStudents)
		sectionRoutes.PUT("/:id", middleware.RequireAdmin(), classHandler.UpdateSection)
		sectionRoutes.DELETE("/:id", middleware.RequireAdmin(), classHandler.DeleteSection)
	}

	// Subjects routes
	subjects := rg.Group("/subjects")
	{
		subjects.GET("", subjectHandler.GetAll)
		subjects.GET("/:id", subjectHandler.GetByID)
		subjects.GET("/class/:classId", subjectHandler.GetByClassID)

		// Admin only routes
		subjects.POST("", middleware.RequireAdmin(), subjectHandler.Create)
		subjects.PUT("/:id", middleware.RequireAdmin(), subjectHandler.Update)
		subjects.DELETE("/:id", middleware.RequireAdmin(), subjectHandler.Delete)
		subjects.POST("/:id/assign-teacher", middleware.RequireAdmin(), subjectHandler.AssignTeacher)
	}

	// Departments routes
	departments := rg.Group("/departments")
	{
		departments.GET("", departmentHandler.GetAll)
		departments.GET("/:id", departmentHandler.GetByID)
		departments.GET("/:id/staff", departmentHandler.GetStaff)

		// Admin only routes
		departments.POST("", middleware.RequireAdmin(), departmentHandler.Create)
		departments.PUT("/:id", middleware.RequireAdmin(), departmentHandler.Update)
		departments.DELETE("/:id", middleware.RequireAdmin(), departmentHandler.Delete)
	}

	// Timetable routes
	timetable := rg.Group("/timetable")
	{
		timetable.GET("", timetableHandler.GetAll)
		timetable.GET("/:id", timetableHandler.GetByID)
		timetable.GET("/class/:classId", timetableHandler.GetByClassID)
		timetable.GET("/section/:sectionId", timetableHandler.GetBySectionID)
		timetable.GET("/teacher/:teacherId", timetableHandler.GetByTeacherID)

		// Admin only routes
		timetable.POST("", middleware.RequireAdmin(), timetableHandler.Create)
		timetable.PUT("/:id", middleware.RequireAdmin(), timetableHandler.Update)
		timetable.DELETE("/:id", middleware.RequireAdmin(), timetableHandler.Delete)
	}
}
