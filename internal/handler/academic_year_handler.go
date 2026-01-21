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

// AcademicYearHandler handles academic year API requests
type AcademicYearHandler struct {
	service *service.AcademicYearService
}

// NewAcademicYearHandler creates a new academic year handler
func NewAcademicYearHandler(service *service.AcademicYearService) *AcademicYearHandler {
	return &AcademicYearHandler{service: service}
}

// Create handles creating a new academic year
func (h *AcademicYearHandler) Create(c *gin.Context) {
	var req request.CreateAcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	resp, err := h.service.Create(&req, institutionID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.Created(c, "Academic year created successfully", resp)
}

// GetAll handles listing all academic years
func (h *AcademicYearHandler) GetAll(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	filter := repository.AcademicYearFilter{
		InstitutionID: middleware.GetInstitutionID(c),
		Search:        c.Query("search"),
	}

	if isCurrent := c.Query("is_current"); isCurrent != "" {
		current := isCurrent == "true"
		filter.IsCurrent = &current
	}

	data, pagination, err := h.service.GetAll(filter, params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

// GetByID handles getting a single academic year
func (h *AcademicYearHandler) GetByID(c *gin.Context) {
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

	resp, err := h.service.GetByID(id, institutionID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}

// GetCurrent handles getting the current academic year
func (h *AcademicYearHandler) GetCurrent(c *gin.Context) {
	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	resp, err := h.service.GetCurrent(institutionID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}

// Update handles updating an academic year
func (h *AcademicYearHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.UpdateAcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	resp, err := h.service.Update(id, &req, institutionID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Academic year updated successfully", resp)
}

// Activate handles setting an academic year as current
func (h *AcademicYearHandler) Activate(c *gin.Context) {
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

	if err := h.service.Activate(id, institutionID); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Academic year activated successfully", nil)
}

// Delete handles deleting an academic year
func (h *AcademicYearHandler) Delete(c *gin.Context) {
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

	if err := h.service.Delete(id, institutionID); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.NoContent(c)
}
