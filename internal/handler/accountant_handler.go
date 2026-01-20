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

// AccountantHandler handles accountant API requests
type AccountantHandler struct {
	service *service.AccountantService
}

func NewAccountantHandler(service *service.AccountantService) *AccountantHandler {
	return &AccountantHandler{service: service}
}

func (h *AccountantHandler) Create(c *gin.Context) {
	var req request.CreateAccountantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	creatorInstID := middleware.GetInstitutionID(c)
	resp, err := h.service.CreateAccountant(&req, creatorInstID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Created(c, "Accountant created successfully", resp)
}

func (h *AccountantHandler) GetAll(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	institutionID := middleware.GetInstitutionID(c)
	data, pagination, err := h.service.GetAllAccountants(institutionID, params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

func (h *AccountantHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	accountant, err := h.service.GetAccountant(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", accountant)
}

func (h *AccountantHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.UpdateAccountantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	institutionID := middleware.GetInstitutionID(c)
	accountant, err := h.service.UpdateAccountant(id, &req, institutionID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Accountant updated successfully", accountant)
}
