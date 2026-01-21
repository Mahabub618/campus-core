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

// TimetableHandler handles timetable API requests
type TimetableHandler struct {
	service *service.TimetableService
}

// NewTimetableHandler creates a new timetable handler
func NewTimetableHandler(service *service.TimetableService) *TimetableHandler {
	return &TimetableHandler{service: service}
}

// Create handles creating a new timetable entry
func (h *TimetableHandler) Create(c *gin.Context) {
	var req request.CreateTimetableRequest
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

	utils.Created(c, "Timetable entry created successfully", resp)
}

// GetAll handles listing all timetable entries
func (h *TimetableHandler) GetAll(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	filter := repository.TimetableFilter{
		InstitutionID:  middleware.GetInstitutionID(c),
		AcademicYearID: c.Query("academic_year_id"),
		ClassID:        c.Query("class_id"),
		SectionID:      c.Query("section_id"),
		SubjectID:      c.Query("subject_id"),
		TeacherID:      c.Query("teacher_id"),
		DayOfWeek:      c.Query("day_of_week"),
	}

	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		filter.IsActive = &active
	}

	data, pagination, err := h.service.GetAll(filter, params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

// GetByID handles getting a single timetable entry
func (h *TimetableHandler) GetByID(c *gin.Context) {
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

// GetByClassID handles getting timetable for a class
func (h *TimetableHandler) GetByClassID(c *gin.Context) {
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

	var academicYearID *uuid.UUID
	if ayIDStr := c.Query("academic_year_id"); ayIDStr != "" {
		ayID, err := uuid.Parse(ayIDStr)
		if err == nil {
			academicYearID = &ayID
		}
	}

	resp, err := h.service.GetByClassID(classID, institutionID, academicYearID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}

// GetBySectionID handles getting timetable for a section
func (h *TimetableHandler) GetBySectionID(c *gin.Context) {
	sectionID, err := uuid.Parse(c.Param("sectionId"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var academicYearID *uuid.UUID
	if ayIDStr := c.Query("academic_year_id"); ayIDStr != "" {
		ayID, err := uuid.Parse(ayIDStr)
		if err == nil {
			academicYearID = &ayID
		}
	}

	resp, err := h.service.GetBySectionID(sectionID, academicYearID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}

// GetByTeacherID handles getting timetable for a teacher
func (h *TimetableHandler) GetByTeacherID(c *gin.Context) {
	teacherID, err := uuid.Parse(c.Param("teacherId"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var academicYearID *uuid.UUID
	if ayIDStr := c.Query("academic_year_id"); ayIDStr != "" {
		ayID, err := uuid.Parse(ayIDStr)
		if err == nil {
			academicYearID = &ayID
		}
	}

	resp, err := h.service.GetByTeacherID(teacherID, academicYearID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", resp)
}

// Update handles updating a timetable entry
func (h *TimetableHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req request.UpdateTimetableRequest
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

	utils.OK(c, "Timetable entry updated successfully", resp)
}

// Delete handles deleting a timetable entry
func (h *TimetableHandler) Delete(c *gin.Context) {
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
