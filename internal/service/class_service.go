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

// ClassService handles class business logic
type ClassService struct {
	classRepo   *repository.ClassRepository
	sectionRepo *repository.SectionRepository
	teacherRepo *repository.TeacherRepository
}

// NewClassService creates a new class service
func NewClassService(classRepo *repository.ClassRepository, sectionRepo *repository.SectionRepository, teacherRepo *repository.TeacherRepository) *ClassService {
	return &ClassService{
		classRepo:   classRepo,
		sectionRepo: sectionRepo,
		teacherRepo: teacherRepo,
	}
}

// CreateClass creates a new class
func (s *ClassService) CreateClass(req *request.CreateClassRequest, institutionID uuid.UUID) (*response.ClassResponse, error) {
	// Check if name already exists
	exists, err := s.classRepo.NameExists(req.Name, institutionID, nil)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}
	if exists {
		return nil, errors.New("class with this name already exists")
	}

	class := &models.Class{
		InstitutionID: institutionID,
		Name:          req.Name,
		Capacity:      req.Capacity,
	}

	// Set class teacher if provided
	if req.ClassTeacherID != "" {
		teacherID, err := uuid.Parse(req.ClassTeacherID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		// Verify teacher exists
		if _, err := s.teacherRepo.FindByID(teacherID); err != nil {
			return nil, errors.New("class teacher not found")
		}
		class.ClassTeacherID = &teacherID
	}

	if err := s.classRepo.Create(class); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return s.toClassResponse(class), nil
}

// GetClassByID gets a class by ID
func (s *ClassService) GetClassByID(id, institutionID uuid.UUID) (*response.ClassResponse, error) {
	class, err := s.classRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return nil, err
	}
	return s.toClassResponse(class), nil
}

// GetAllClasses gets all classes with filters
func (s *ClassService) GetAllClasses(filter repository.ClassFilter, params utils.PaginationParams) ([]response.ClassResponse, utils.Pagination, error) {
	classes, total, err := s.classRepo.FindAll(filter, params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.ClassResponse
	for _, class := range classes {
		responses = append(responses, *s.toClassResponse(&class))
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return responses, pagination, nil
}

// UpdateClass updates a class
func (s *ClassService) UpdateClass(id uuid.UUID, req *request.UpdateClassRequest, institutionID uuid.UUID) (*response.ClassResponse, error) {
	class, err := s.classRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" && req.Name != class.Name {
		exists, err := s.classRepo.NameExists(req.Name, institutionID, &id)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, errors.New("class with this name already exists")
		}
		class.Name = req.Name
	}

	if req.Capacity != nil {
		class.Capacity = *req.Capacity
	}

	if req.ClassTeacherID != "" {
		teacherID, err := uuid.Parse(req.ClassTeacherID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		// Verify teacher exists
		if _, err := s.teacherRepo.FindByID(teacherID); err != nil {
			return nil, errors.New("class teacher not found")
		}
		class.ClassTeacherID = &teacherID
	}

	if err := s.classRepo.Update(class); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return s.toClassResponse(class), nil
}

// DeleteClass deletes a class
func (s *ClassService) DeleteClass(id, institutionID uuid.UUID) error {
	// Verify it exists and belongs to the institution
	_, err := s.classRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return err
	}

	// Check if class has students
	count, err := s.classRepo.GetClassStudentCount(id)
	if err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}
	if count > 0 {
		return errors.New("cannot delete class with students")
	}

	return s.classRepo.Delete(id)
}

// GetClassStudents gets all students in a class
func (s *ClassService) GetClassStudents(classID, institutionID uuid.UUID) ([]response.UserResponse, error) {
	// Verify class exists and belongs to the institution
	_, err := s.classRepo.FindByIDWithInstitution(classID, institutionID)
	if err != nil {
		return nil, err
	}

	// This would need a student service/repository integration
	// For now, return empty slice - will be implemented with student service
	return []response.UserResponse{}, nil
}

// GetClassTeachers gets all teachers assigned to a class
func (s *ClassService) GetClassTeachers(classID, institutionID uuid.UUID) ([]response.TeacherBrief, error) {
	// Verify class exists and belongs to the institution
	_, err := s.classRepo.FindByIDWithInstitution(classID, institutionID)
	if err != nil {
		return nil, err
	}

	teachers, err := s.classRepo.GetClassTeachers(classID)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.TeacherBrief
	for _, t := range teachers {
		responses = append(responses, s.toTeacherBrief(&t))
	}

	return responses, nil
}

// CreateSection creates a new section for a class
func (s *ClassService) CreateSection(classID uuid.UUID, req *request.CreateSectionRequest, institutionID uuid.UUID) (*response.SectionResponse, error) {
	// Verify class exists and belongs to the institution
	class, err := s.classRepo.FindByIDWithInstitution(classID, institutionID)
	if err != nil {
		return nil, err
	}

	// Check if section name already exists in class
	exists, err := s.sectionRepo.NameExistsInClass(req.Name, classID, nil)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}
	if exists {
		return nil, errors.New("section with this name already exists in class")
	}

	section := &models.Section{
		ClassID:    classID,
		Name:       req.Name,
		RoomNumber: req.RoomNumber,
		Capacity:   req.Capacity,
	}

	if err := s.sectionRepo.Create(section); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// Update section count in class
	class.SectionCount++
	_ = s.classRepo.Update(class)

	return s.toSectionResponse(section), nil
}

// GetSectionsByClass gets all sections for a class
func (s *ClassService) GetSectionsByClass(classID, institutionID uuid.UUID) ([]response.SectionResponse, error) {
	// Verify class exists and belongs to the institution
	_, err := s.classRepo.FindByIDWithInstitution(classID, institutionID)
	if err != nil {
		return nil, err
	}

	sections, err := s.sectionRepo.FindByClassID(classID)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.SectionResponse
	for _, section := range sections {
		responses = append(responses, *s.toSectionResponse(&section))
	}

	return responses, nil
}

// UpdateSection updates a section
func (s *ClassService) UpdateSection(sectionID uuid.UUID, req *request.UpdateSectionRequest) (*response.SectionResponse, error) {
	section, err := s.sectionRepo.FindByID(sectionID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" && req.Name != section.Name {
		exists, err := s.sectionRepo.NameExistsInClass(req.Name, section.ClassID, &sectionID)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, errors.New("section with this name already exists in class")
		}
		section.Name = req.Name
	}

	if req.RoomNumber != "" {
		section.RoomNumber = req.RoomNumber
	}
	if req.Capacity != nil {
		section.Capacity = *req.Capacity
	}

	if err := s.sectionRepo.Update(section); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return s.toSectionResponse(section), nil
}

// DeleteSection deletes a section
func (s *ClassService) DeleteSection(sectionID uuid.UUID) error {
	section, err := s.sectionRepo.FindByID(sectionID)
	if err != nil {
		return err
	}

	// Check if section has students
	count, err := s.sectionRepo.GetSectionStudentCount(sectionID)
	if err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}
	if count > 0 {
		return errors.New("cannot delete section with students")
	}

	// Update section count in class
	class, err := s.classRepo.FindByID(section.ClassID)
	if err == nil && class.SectionCount > 0 {
		class.SectionCount--
		_ = s.classRepo.Update(class)
	}

	return s.sectionRepo.Delete(sectionID)
}

// GetSectionStudents gets all students in a section
func (s *ClassService) GetSectionStudents(sectionID uuid.UUID) ([]response.UserResponse, error) {
	// Verify section exists
	_, err := s.sectionRepo.FindByID(sectionID)
	if err != nil {
		return nil, err
	}

	// This would need a student service/repository integration
	// For now, return empty slice - will be implemented with student service
	return []response.UserResponse{}, nil
}

// Helper methods for converting models to responses
func (s *ClassService) toClassResponse(class *models.Class) *response.ClassResponse {
	resp := &response.ClassResponse{
		ID:            class.ID,
		InstitutionID: class.InstitutionID,
		Name:          class.Name,
		SectionCount:  class.SectionCount,
		Capacity:      class.Capacity,
		CreatedAt:     class.CreatedAt,
		UpdatedAt:     class.UpdatedAt,
	}

	if class.ClassTeacherID != nil {
		resp.ClassTeacherID = class.ClassTeacherID
		if class.ClassTeacher != nil {
			brief := s.toTeacherBrief(class.ClassTeacher)
			resp.ClassTeacher = &brief
		}
	}

	if len(class.Sections) > 0 {
		for _, section := range class.Sections {
			resp.Sections = append(resp.Sections, *s.toSectionResponse(&section))
		}
	}

	return resp
}

func (s *ClassService) toSectionResponse(section *models.Section) *response.SectionResponse {
	resp := &response.SectionResponse{
		ID:         section.ID,
		ClassID:    section.ClassID,
		Name:       section.Name,
		RoomNumber: section.RoomNumber,
		Capacity:   section.Capacity,
		CreatedAt:  section.CreatedAt,
		UpdatedAt:  section.UpdatedAt,
	}

	if section.Class != nil {
		resp.Class = &response.ClassBrief{
			ID:   section.Class.ID,
			Name: section.Class.Name,
		}
	}

	return resp
}

func (s *ClassService) toTeacherBrief(teacher *models.Teacher) response.TeacherBrief {
	brief := response.TeacherBrief{
		ID: teacher.ID,
	}
	if teacher.User != nil && teacher.User.Profile != nil {
		brief.FirstName = teacher.User.Profile.FirstName
		brief.LastName = teacher.User.Profile.LastName
	}
	return brief
}
