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

// TeacherHandler handles teacher API requests
type TeacherHandler struct {
	service *service.TeacherService
}

func NewTeacherHandler(service *service.TeacherService) *TeacherHandler {
	return &TeacherHandler{service: service}
}

func (h *TeacherHandler) Create(c *gin.Context) {
	var req request.CreateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	creatorInstID := middleware.GetInstitutionID(c)
	resp, err := h.service.CreateTeacher(&req, creatorInstID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err) // Can be duplicate key error etc.
		return
	}

	utils.Created(c, "Teacher created successfully", resp)
}

func (h *TeacherHandler) GetAll(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	institutionID := middleware.GetInstitutionID(c)
	data, pagination, err := h.service.GetAllTeachers(institutionID, params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

func (h *TeacherHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	teacher, err := h.service.GetTeacher(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", teacher)
}
