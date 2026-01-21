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

// SubjectHandler handles subject API requests
type SubjectHandler struct {
	service *service.SubjectService
}

// NewSubjectHandler creates a new subject handler
func NewSubjectHandler(service *service.SubjectService) *SubjectHandler {
	return &SubjectHandler{service: service}
}

// Create handles creating a new subject
func (h *SubjectHandler) Create(c *gin.Context) {
	var req request.CreateSubjectRequest
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

	utils.Created(c, "Subject created successfully", resp)
}

// GetAll handles listing all subjects
func (h *SubjectHandler) GetAll(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	filter := repository.SubjectFilter{
		InstitutionID: middleware.GetInstitutionID(c),
		ClassID:       c.Query("class_id"),
		TeacherID:     c.Query("teacher_id"),
		Search:        c.Query("search"),
	}

	if isElective := c.Query("is_elective"); isElective != "" {
		elective := isElective == "true"
		filter.IsElective = &elective
	}

	data, pagination, err := h.service.GetAll(filter, params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

// GetByID handles getting a single subject
func (h *SubjectHandler) GetByID(c *gin.Context) {
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

// GetByClassID handles getting subjects for a specific class
func (h *SubjectHandler) GetByClassID(c *gin.Context) {
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

	resp, err := h.service.GetByClassID(classID, institutionID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}

// Update handles updating a subject
func (h *SubjectHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.UpdateSubjectRequest
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

	utils.OK(c, "Subject updated successfully", resp)
}

// Delete handles deleting a subject
func (h *SubjectHandler) Delete(c *gin.Context) {
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

// AssignTeacher handles assigning a teacher to a subject
func (h *SubjectHandler) AssignTeacher(c *gin.Context) {
	subjectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.AssignTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	institutionID, err := uuid.Parse(middleware.GetInstitutionID(c))
	if err != nil {
		utils.BadRequest(c, "Invalid institution ID")
		return
	}

	if err := h.service.AssignTeacher(subjectID, &req, institutionID); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Teacher assigned successfully", nil)
}
