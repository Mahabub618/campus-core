package service

import (
	"errors"

	"campus-core/internal/dto/request"
	"campus-core/internal/dto/response"
	"campus-core/internal/models"
	"campus-core/internal/repository"
	"campus-core/internal/utils"

	"github.com/google/uuid"
)

// TimetableService handles timetable business logic
type TimetableService struct {
	ttRepo      *repository.TimetableRepository
	classRepo   *repository.ClassRepository
	sectionRepo *repository.SectionRepository
	subjectRepo *repository.SubjectRepository
	teacherRepo *repository.TeacherRepository
	ayRepo      *repository.AcademicYearRepository
}

// NewTimetableService creates a new timetable service
func NewTimetableService(
	ttRepo *repository.TimetableRepository,
	classRepo *repository.ClassRepository,
	sectionRepo *repository.SectionRepository,
	subjectRepo *repository.SubjectRepository,
	teacherRepo *repository.TeacherRepository,
	ayRepo *repository.AcademicYearRepository,
) *TimetableService {
	return &TimetableService{
		ttRepo:      ttRepo,
		classRepo:   classRepo,
		sectionRepo: sectionRepo,
		subjectRepo: subjectRepo,
		teacherRepo: teacherRepo,
		ayRepo:      ayRepo,
	}
}

// Create creates a new timetable entry
func (s *TimetableService) Create(req *request.CreateTimetableRequest, institutionID uuid.UUID) (*response.TimetableResponse, error) {
	// Parse and validate all UUIDs
	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		return nil, utils.ErrInvalidUUID
	}
	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		return nil, utils.ErrInvalidUUID
	}
	sectionID, err := uuid.Parse(req.SectionID)
	if err != nil {
		return nil, utils.ErrInvalidUUID
	}
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, utils.ErrInvalidUUID
	}
	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		return nil, utils.ErrInvalidUUID
	}

	// Verify all entities exist
	if _, err := s.ayRepo.FindByIDWithInstitution(academicYearID, institutionID); err != nil {
		return nil, errors.New("academic year not found")
	}
	if _, err := s.classRepo.FindByIDWithInstitution(classID, institutionID); err != nil {
		return nil, errors.New("class not found")
	}
	if _, err := s.sectionRepo.FindByID(sectionID); err != nil {
		return nil, errors.New("section not found")
	}
	if _, err := s.subjectRepo.FindByIDWithInstitution(subjectID, institutionID); err != nil {
		return nil, errors.New("subject not found")
	}
	if _, err := s.teacherRepo.FindByID(teacherID); err != nil {
		return nil, errors.New("teacher not found")
	}

	tt := &models.Timetable{
		InstitutionID:  institutionID,
		AcademicYearID: academicYearID,
		ClassID:        classID,
		SectionID:      sectionID,
		SubjectID:      subjectID,
		TeacherID:      teacherID,
		DayOfWeek:      models.DayOfWeek(req.DayOfWeek),
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		RoomNumber:     req.RoomNumber,
		IsActive:       true,
	}

	// Check for conflicts
	hasConflict, err := s.ttRepo.CheckConflict(tt, nil)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}
	if hasConflict {
		return nil, errors.New("scheduling conflict detected: teacher, section, or room is already occupied at this time")
	}

	if err := s.ttRepo.Create(tt); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// Reload with preloads
	tt, _ = s.ttRepo.FindByID(tt.ID)

	return s.toResponse(tt), nil
}

// GetByID gets a timetable entry by ID
func (s *TimetableService) GetByID(id, institutionID uuid.UUID) (*response.TimetableResponse, error) {
	tt, err := s.ttRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return nil, err
	}
	return s.toResponse(tt), nil
}

// GetAll gets all timetable entries with filters
func (s *TimetableService) GetAll(filter repository.TimetableFilter, params utils.PaginationParams) ([]response.TimetableResponse, utils.Pagination, error) {
	timetables, total, err := s.ttRepo.FindAll(filter, params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.TimetableResponse
	for _, tt := range timetables {
		responses = append(responses, *s.toResponse(&tt))
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return responses, pagination, nil
}

// GetByClassID gets timetable for a class
func (s *TimetableService) GetByClassID(classID, institutionID uuid.UUID, academicYearID *uuid.UUID) (*response.WeekTimetableResponse, error) {
	// Verify class exists
	if _, err := s.classRepo.FindByIDWithInstitution(classID, institutionID); err != nil {
		return nil, err
	}

	timetables, err := s.ttRepo.FindByClassID(classID, academicYearID)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return s.groupByDay(timetables), nil
}

// GetBySectionID gets timetable for a section
func (s *TimetableService) GetBySectionID(sectionID uuid.UUID, academicYearID *uuid.UUID) (*response.WeekTimetableResponse, error) {
	// Verify section exists
	if _, err := s.sectionRepo.FindByID(sectionID); err != nil {
		return nil, err
	}

	timetables, err := s.ttRepo.FindBySectionID(sectionID, academicYearID)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return s.groupByDay(timetables), nil
}

// GetByTeacherID gets timetable for a teacher
func (s *TimetableService) GetByTeacherID(teacherID uuid.UUID, academicYearID *uuid.UUID) (*response.WeekTimetableResponse, error) {
	// Verify teacher exists
	if _, err := s.teacherRepo.FindByID(teacherID); err != nil {
		return nil, err
	}

	timetables, err := s.ttRepo.FindByTeacherID(teacherID, academicYearID)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return s.groupByDay(timetables), nil
}

// Update updates a timetable entry
func (s *TimetableService) Update(id uuid.UUID, req *request.UpdateTimetableRequest, institutionID uuid.UUID) (*response.TimetableResponse, error) {
	tt, err := s.ttRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.AcademicYearID != "" {
		ayID, err := uuid.Parse(req.AcademicYearID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		if _, err := s.ayRepo.FindByIDWithInstitution(ayID, institutionID); err != nil {
			return nil, errors.New("academic year not found")
		}
		tt.AcademicYearID = ayID
	}
	if req.ClassID != "" {
		classID, err := uuid.Parse(req.ClassID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		if _, err := s.classRepo.FindByIDWithInstitution(classID, institutionID); err != nil {
			return nil, errors.New("class not found")
		}
		tt.ClassID = classID
	}
	if req.SectionID != "" {
		sectionID, err := uuid.Parse(req.SectionID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		if _, err := s.sectionRepo.FindByID(sectionID); err != nil {
			return nil, errors.New("section not found")
		}
		tt.SectionID = sectionID
	}
	if req.SubjectID != "" {
		subjectID, err := uuid.Parse(req.SubjectID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		if _, err := s.subjectRepo.FindByIDWithInstitution(subjectID, institutionID); err != nil {
			return nil, errors.New("subject not found")
		}
		tt.SubjectID = subjectID
	}
	if req.TeacherID != "" {
		teacherID, err := uuid.Parse(req.TeacherID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		if _, err := s.teacherRepo.FindByID(teacherID); err != nil {
			return nil, errors.New("teacher not found")
		}
		tt.TeacherID = teacherID
	}
	if req.DayOfWeek != "" {
		tt.DayOfWeek = models.DayOfWeek(req.DayOfWeek)
	}
	if req.StartTime != "" {
		tt.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		tt.EndTime = req.EndTime
	}
	if req.RoomNumber != "" {
		tt.RoomNumber = req.RoomNumber
	}
	if req.IsActive != nil {
		tt.IsActive = *req.IsActive
	}

	// Check for conflicts
	hasConflict, err := s.ttRepo.CheckConflict(tt, &id)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}
	if hasConflict {
		return nil, errors.New("scheduling conflict detected: teacher, section, or room is already occupied at this time")
	}

	if err := s.ttRepo.Update(tt); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// Reload with preloads
	tt, _ = s.ttRepo.FindByID(tt.ID)

	return s.toResponse(tt), nil
}

// Delete deletes a timetable entry
func (s *TimetableService) Delete(id, institutionID uuid.UUID) error {
	// Verify it exists and belongs to the institution
	_, err := s.ttRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return err
	}

	return s.ttRepo.Delete(id)
}

// groupByDay groups timetable entries by day of week
func (s *TimetableService) groupByDay(timetables []models.Timetable) *response.WeekTimetableResponse {
	dayOrder := []string{"SUNDAY", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY"}
	dayMap := make(map[string][]response.TimetableResponse)

	for _, tt := range timetables {
		day := string(tt.DayOfWeek)
		dayMap[day] = append(dayMap[day], *s.toResponse(&tt))
	}

	var days []response.DayTimetable
	for _, day := range dayOrder {
		if entries, ok := dayMap[day]; ok {
			days = append(days, response.DayTimetable{
				Day:     day,
				Entries: entries,
			})
		}
	}

	return &response.WeekTimetableResponse{Days: days}
}

// toResponse converts a model to response
func (s *TimetableService) toResponse(tt *models.Timetable) *response.TimetableResponse {
	resp := &response.TimetableResponse{
		ID:             tt.ID,
		InstitutionID:  tt.InstitutionID,
		AcademicYearID: tt.AcademicYearID,
		ClassID:        tt.ClassID,
		SectionID:      tt.SectionID,
		SubjectID:      tt.SubjectID,
		TeacherID:      tt.TeacherID,
		DayOfWeek:      string(tt.DayOfWeek),
		StartTime:      tt.StartTime,
		EndTime:        tt.EndTime,
		RoomNumber:     tt.RoomNumber,
		IsActive:       tt.IsActive,
		CreatedAt:      tt.CreatedAt,
		UpdatedAt:      tt.UpdatedAt,
	}

	if tt.Class != nil {
		resp.Class = &response.ClassBrief{
			ID:   tt.Class.ID,
			Name: tt.Class.Name,
		}
	}
	if tt.Section != nil {
		resp.Section = &response.SectionBrief{
			ID:   tt.Section.ID,
			Name: tt.Section.Name,
		}
	}
	if tt.Subject != nil {
		resp.Subject = &response.SubjectBrief{
			ID:   tt.Subject.ID,
			Name: tt.Subject.Name,
			Code: tt.Subject.Code,
		}
	}
	if tt.Teacher != nil {
		resp.Teacher = &response.TeacherBrief{
			ID: tt.Teacher.ID,
		}
		if tt.Teacher.User != nil && tt.Teacher.User.Profile != nil {
			resp.Teacher.FirstName = tt.Teacher.User.Profile.FirstName
			resp.Teacher.LastName = tt.Teacher.User.Profile.LastName
		}
	}

	return resp
}
