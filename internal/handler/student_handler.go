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

// StudentHandler handles student API requests
type StudentHandler struct {
	service *service.StudentService
}

func NewStudentHandler(service *service.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req request.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	creatorInstID := middleware.GetInstitutionID(c)
	resp, err := h.service.CreateStudent(&req, creatorInstID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Created(c, "Student created successfully", resp)
}

func (h *StudentHandler) GetAll(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	institutionID := middleware.GetInstitutionID(c)
	data, pagination, err := h.service.GetAllStudents(institutionID, params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

func (h *StudentHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	student, err := h.service.GetStudent(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", student)
}

func (h *StudentHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	institutionID := middleware.GetInstitutionID(c)
	student, err := h.service.UpdateStudent(id, &req, institutionID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Student updated successfully", student)
}

func (h *StudentHandler) GetParents(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	parents, err := h.service.GetStudentParents(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", parents)
}

func (h *StudentHandler) LinkParent(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.LinkParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.LinkParent(studentID, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.Created(c, "Parent linked successfully", nil)
}

func (h *StudentHandler) UnlinkParent(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	parentID, err := uuid.Parse(c.Param("parentId"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	if err := h.service.UnlinkParent(studentID, parentID); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Parent unlinked successfully", nil)
}
