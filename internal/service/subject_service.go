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

// SubjectService handles subject business logic
type SubjectService struct {
	subjectRepo *repository.SubjectRepository
	classRepo   *repository.ClassRepository
	teacherRepo *repository.TeacherRepository
}

// NewSubjectService creates a new subject service
func NewSubjectService(subjectRepo *repository.SubjectRepository, classRepo *repository.ClassRepository, teacherRepo *repository.TeacherRepository) *SubjectService {
	return &SubjectService{
		subjectRepo: subjectRepo,
		classRepo:   classRepo,
		teacherRepo: teacherRepo,
	}
}

// Create creates a new subject
func (s *SubjectService) Create(req *request.CreateSubjectRequest, institutionID uuid.UUID) (*response.SubjectResponse, error) {
	subject := &models.Subject{
		InstitutionID: institutionID,
		Name:          req.Name,
		Code:          req.Code,
		IsElective:    req.IsElective,
		CreditHours:   req.CreditHours,
	}

	// Set class if provided
	if req.ClassID != "" {
		classID, err := uuid.Parse(req.ClassID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		// Verify class exists and belongs to institution
		if _, err := s.classRepo.FindByIDWithInstitution(classID, institutionID); err != nil {
			return nil, errors.New("class not found")
		}
		subject.ClassID = &classID

		// Check if subject name already exists in class
		exists, err := s.subjectRepo.NameExistsInClass(req.Name, classID, nil)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, errors.New("subject with this name already exists in class")
		}
	}

	// Set teacher if provided
	if req.TeacherID != "" {
		teacherID, err := uuid.Parse(req.TeacherID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		// Verify teacher exists
		if _, err := s.teacherRepo.FindByID(teacherID); err != nil {
			return nil, errors.New("teacher not found")
		}
		subject.TeacherID = &teacherID
	}

	// Check if code already exists (if provided)
	if req.Code != "" {
		exists, err := s.subjectRepo.CodeExists(req.Code, institutionID, nil)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, errors.New("subject with this code already exists")
		}
	}

	if err := s.subjectRepo.Create(subject); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return s.toResponse(subject), nil
}

// GetByID gets a subject by ID
func (s *SubjectService) GetByID(id, institutionID uuid.UUID) (*response.SubjectResponse, error) {
	subject, err := s.subjectRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return nil, err
	}
	return s.toResponse(subject), nil
}

// GetAll gets all subjects with filters
func (s *SubjectService) GetAll(filter repository.SubjectFilter, params utils.PaginationParams) ([]response.SubjectResponse, utils.Pagination, error) {
	subjects, total, err := s.subjectRepo.FindAll(filter, params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.SubjectResponse
	for _, subject := range subjects {
		responses = append(responses, *s.toResponse(&subject))
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return responses, pagination, nil
}

// GetByClassID gets all subjects for a class
func (s *SubjectService) GetByClassID(classID, institutionID uuid.UUID) ([]response.SubjectResponse, error) {
	// Verify class exists and belongs to institution
	if _, err := s.classRepo.FindByIDWithInstitution(classID, institutionID); err != nil {
		return nil, err
	}

	subjects, err := s.subjectRepo.FindByClassID(classID)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.SubjectResponse
	for _, subject := range subjects {
		responses = append(responses, *s.toResponse(&subject))
	}

	return responses, nil
}

// Update updates a subject
func (s *SubjectService) Update(id uuid.UUID, req *request.UpdateSubjectRequest, institutionID uuid.UUID) (*response.SubjectResponse, error) {
	subject, err := s.subjectRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return nil, err
	}

	// Update name if provided
	if req.Name != "" && req.Name != subject.Name {
		if subject.ClassID != nil {
			exists, err := s.subjectRepo.NameExistsInClass(req.Name, *subject.ClassID, &id)
			if err != nil {
				return nil, utils.ErrInternalServer.Wrap(err)
			}
			if exists {
				return nil, errors.New("subject with this name already exists in class")
			}
		}
		subject.Name = req.Name
	}

	// Update code if provided
	if req.Code != "" && req.Code != subject.Code {
		exists, err := s.subjectRepo.CodeExists(req.Code, institutionID, &id)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, errors.New("subject with this code already exists")
		}
		subject.Code = req.Code
	}

	// Update class if provided
	if req.ClassID != "" {
		classID, err := uuid.Parse(req.ClassID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		if _, err := s.classRepo.FindByIDWithInstitution(classID, institutionID); err != nil {
			return nil, errors.New("class not found")
		}
		subject.ClassID = &classID
	}

	// Update teacher if provided
	if req.TeacherID != "" {
		teacherID, err := uuid.Parse(req.TeacherID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		if _, err := s.teacherRepo.FindByID(teacherID); err != nil {
			return nil, errors.New("teacher not found")
		}
		subject.TeacherID = &teacherID
	}

	if req.IsElective != nil {
		subject.IsElective = *req.IsElective
	}
	if req.CreditHours != nil {
		subject.CreditHours = *req.CreditHours
	}

	if err := s.subjectRepo.Update(subject); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return s.toResponse(subject), nil
}

// Delete deletes a subject
func (s *SubjectService) Delete(id, institutionID uuid.UUID) error {
	// Verify it exists and belongs to the institution
	_, err := s.subjectRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return err
	}

	return s.subjectRepo.Delete(id)
}

// AssignTeacher assigns a teacher to a subject
func (s *SubjectService) AssignTeacher(subjectID uuid.UUID, req *request.AssignTeacherRequest, institutionID uuid.UUID) error {
	// Verify subject exists and belongs to the institution
	_, err := s.subjectRepo.FindByIDWithInstitution(subjectID, institutionID)
	if err != nil {
		return err
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		return utils.ErrInvalidUUID
	}

	// Verify teacher exists
	if _, err := s.teacherRepo.FindByID(teacherID); err != nil {
		return errors.New("teacher not found")
	}

	return s.subjectRepo.AssignTeacher(subjectID, teacherID)
}

// toResponse converts a model to response
func (s *SubjectService) toResponse(subject *models.Subject) *response.SubjectResponse {
	resp := &response.SubjectResponse{
		ID:            subject.ID,
		InstitutionID: subject.InstitutionID,
		Name:          subject.Name,
		Code:          subject.Code,
		IsElective:    subject.IsElective,
		CreditHours:   subject.CreditHours,
		CreatedAt:     subject.CreatedAt,
		UpdatedAt:     subject.UpdatedAt,
	}

	if subject.ClassID != nil {
		resp.ClassID = subject.ClassID
		if subject.Class != nil {
			resp.Class = &response.ClassBrief{
				ID:   subject.Class.ID,
				Name: subject.Class.Name,
			}
		}
	}

	if subject.TeacherID != nil {
		resp.TeacherID = subject.TeacherID
		if subject.Teacher != nil {
			resp.Teacher = &response.TeacherBrief{
				ID: subject.Teacher.ID,
			}
			if subject.Teacher.User != nil && subject.Teacher.User.Profile != nil {
				resp.Teacher.FirstName = subject.Teacher.User.Profile.FirstName
				resp.Teacher.LastName = subject.Teacher.User.Profile.LastName
			}
		}
	}

	return resp
}
