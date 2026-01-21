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

// DepartmentHandler handles department API requests
type DepartmentHandler struct {
	service *service.DepartmentService
}

// NewDepartmentHandler creates a new department handler
func NewDepartmentHandler(service *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// Create handles creating a new department
func (h *DepartmentHandler) Create(c *gin.Context) {
	var req request.CreateDepartmentRequest
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

	utils.Created(c, "Department created successfully", resp)
}

// GetAll handles listing all departments
func (h *DepartmentHandler) GetAll(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	filter := repository.DepartmentFilter{
		InstitutionID: middleware.GetInstitutionID(c),
		Search:        c.Query("search"),
	}

	data, pagination, err := h.service.GetAll(filter, params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

// GetByID handles getting a single department
func (h *DepartmentHandler) GetByID(c *gin.Context) {
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

// Update handles updating a department
func (h *DepartmentHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.UpdateDepartmentRequest
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

	utils.OK(c, "Department updated successfully", resp)
}

// Delete handles deleting a department
func (h *DepartmentHandler) Delete(c *gin.Context) {
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

// GetStaff handles getting all staff in a department
func (h *DepartmentHandler) GetStaff(c *gin.Context) {
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

	resp, err := h.service.GetDepartmentStaff(id, institutionID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}
