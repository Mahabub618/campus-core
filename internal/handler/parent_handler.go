package handler

import (
	"net/http"

	"campus-core/internal/dto/request"
	"campus-core/internal/middleware"
	"campus-core/internal/service"
	"campus-core/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ParentHandler handles parent API requests
type ParentHandler struct {
	service *service.ParentService
}

func NewParentHandler(service *service.ParentService) *ParentHandler {
	return &ParentHandler{service: service}
}

func (h *ParentHandler) Create(c *gin.Context) {
	var req request.CreateParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	creatorInstID := middleware.GetInstitutionID(c)
	resp, err := h.service.CreateParent(&req, creatorInstID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Created(c, "Parent created successfully", resp)
}

func (h *ParentHandler) GetAll(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	institutionID := middleware.GetInstitutionID(c)
	data, pagination, err := h.service.GetAllParents(institutionID, params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

func (h *ParentHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	parent, err := h.service.GetParent(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", parent)
}

func (h *ParentHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.UpdateParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	institutionID := middleware.GetInstitutionID(c)
	parent, err := h.service.UpdateParent(id, &req, institutionID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Parent updated successfully", parent)
}

func (h *ParentHandler) GetChildren(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	children, err := h.service.GetParentChildren(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", children)
}
