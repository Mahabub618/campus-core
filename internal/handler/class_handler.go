package handler

import (
	"net/http"

	"campus-core/internal/dto/request"
	"campus-core/internal/middleware"
	"campus-core/internal/repository"
	"campus-core/internal/service"
	"campus-core/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ClassHandler handles class API requests
type ClassHandler struct {
	service *service.ClassService
}

// NewClassHandler creates a new class handler
func NewClassHandler(service *service.ClassService) *ClassHandler {
	return &ClassHandler{service: service}
}

// Create handles creating a new class
func (h *ClassHandler) Create(c *gin.Context) {
	var req request.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	resp, err := h.service.CreateClass(&req, institutionID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.Created(c, "Class created successfully", resp)
}

// GetAll handles listing all classes
func (h *ClassHandler) GetAll(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	filter := repository.ClassFilter{
		InstitutionID: middleware.GetInstitutionID(c),
		Search:        c.Query("search"),
	}

	data, pagination, err := h.service.GetAllClasses(filter, params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

// GetByID handles getting a single class
func (h *ClassHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	resp, err := h.service.GetClassByID(id, institutionID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}

// Update handles updating a class
func (h *ClassHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	resp, err := h.service.UpdateClass(id, &req, institutionID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Class updated successfully", resp)
}

// Delete handles deleting a class
func (h *ClassHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	if err := h.service.DeleteClass(id, institutionID); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.NoContent(c)
}

// GetStudents handles getting all students in a class
func (h *ClassHandler) GetStudents(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	resp, err := h.service.GetClassStudents(id, institutionID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}

// GetTeachers handles getting all teachers for a class
func (h *ClassHandler) GetTeachers(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	resp, err := h.service.GetClassTeachers(id, institutionID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}

// CreateSection handles creating a new section for a class
func (h *ClassHandler) CreateSection(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("classId"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.CreateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	resp, err := h.service.CreateSection(classID, &req, institutionID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.Created(c, "Section created successfully", resp)
}

// GetSections handles getting all sections for a class
func (h *ClassHandler) GetSections(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("classId"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	resp, err := h.service.GetSectionsByClass(classID, institutionID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}

// UpdateSection handles updating a section
func (h *ClassHandler) UpdateSection(c *gin.Context) {
	sectionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.UpdateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.UpdateSection(sectionID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Section updated successfully", resp)
}

// DeleteSection handles deleting a section
func (h *ClassHandler) DeleteSection(c *gin.Context) {
	sectionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	if err := h.service.DeleteSection(sectionID); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.NoContent(c)
}

// GetSectionStudents handles getting all students in a section
func (h *ClassHandler) GetSectionStudents(c *gin.Context) {
	sectionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	resp, err := h.service.GetSectionStudents(sectionID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}
