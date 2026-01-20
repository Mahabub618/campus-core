package handler

import (
	"net/http"

	"campus-core/internal/models"
	"campus-core/internal/service"
	"campus-core/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InstitutionHandler handles institution API requests
type InstitutionHandler struct {
	service *service.InstitutionService
}

// NewInstitutionHandler creates a new institution handler
func NewInstitutionHandler(service *service.InstitutionService) *InstitutionHandler {
	return &InstitutionHandler{service: service}
}

// Create handles institution creation
func (h *InstitutionHandler) Create(c *gin.Context) {
	var input struct {
		Name          string `json:"name" binding:"required"`
		Code          string `json:"code" binding:"required"`
		Address       string `json:"address"`
		Phone         string `json:"phone"`
		Email         string `json:"email" binding:"omitempty,email"`
		PrincipalName string `json:"principal_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	institution := &models.Institution{
		Name:          input.Name,
		Code:          input.Code,
		Address:       input.Address,
		Phone:         input.Phone,
		Email:         input.Email,
		PrincipalName: input.PrincipalName,
		IsActive:      true,
	}

	if err := h.service.Create(institution); err != nil {
		utils.Error(c, http.StatusBadRequest, err) // Could be 409 if code exists
		return
	}

	utils.Created(c, "Institution created successfully", institution)
}

// GetAll returns all institutions
func (h *InstitutionHandler) GetAll(c *gin.Context) {
	var params utils.PaginationParams
	// Manual binding for query params as ShouldBindQuery might be overkill or we use our helper
	// Assuming default binding works for int
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	data, pagination, err := h.service.GetAll(params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

// GetByID returns an institution by ID
func (h *InstitutionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	institution, err := h.service.GetByID(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", institution)
}

// Update updates an institution
func (h *InstitutionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	institution, err := h.service.Update(id, updates)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Institution updated successfully", institution)
}

// Delete deletes an institution
func (h *InstitutionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.OK(c, "Institution deleted successfully", nil)
}

// GetStats returns institution stats
func (h *InstitutionHandler) GetStats(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	stats, err := h.service.GetStats(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.OK(c, "", stats)
}

// ToggleStatus enables or disables an institution
func (h *InstitutionHandler) ToggleStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.ToggleStatus(id, req.IsActive); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	status := "disabled"
	if req.IsActive {
		status = "enabled"
	}
	utils.OK(c, "Institution "+status+" successfully", nil)
}

// GetAdmins returns a list of admins for an institution
func (h *InstitutionHandler) GetAdmins(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	admins, err := h.service.GetAdmins(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.OK(c, "", admins)
}

// AssignAdmin assigns an admin to an institution
func (h *InstitutionHandler) AssignAdmin(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req struct {
		Email     string `json:"email" binding:"required,email"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Password  string `json:"password" binding:"required,min=8"`
		Phone     string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	admin, err := h.service.AssignAdmin(id, req.Email, req.FirstName, req.LastName, req.Password, req.Phone)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.Created(c, "Admin assigned successfully", admin)
}
