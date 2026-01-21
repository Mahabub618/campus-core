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

// DepartmentService handles department business logic
type DepartmentService struct {
	deptRepo    *repository.DepartmentRepository
	teacherRepo *repository.TeacherRepository
}

// NewDepartmentService creates a new department service
func NewDepartmentService(deptRepo *repository.DepartmentRepository, teacherRepo *repository.TeacherRepository) *DepartmentService {
	return &DepartmentService{
		deptRepo:    deptRepo,
		teacherRepo: teacherRepo,
	}
}

// Create creates a new department
func (s *DepartmentService) Create(req *request.CreateDepartmentRequest, institutionID uuid.UUID) (*response.DepartmentResponse, error) {
	// Check if name already exists
	exists, err := s.deptRepo.NameExists(req.Name, institutionID, nil)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}
	if exists {
		return nil, errors.New("department with this name already exists")
	}

	dept := &models.Department{
		InstitutionID: institutionID,
		Name:          req.Name,
		Description:   req.Description,
	}

	// Set head of department if provided
	if req.HeadOfDepartmentID != "" {
		hodID, err := uuid.Parse(req.HeadOfDepartmentID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		// Verify teacher exists
		if _, err := s.teacherRepo.FindByID(hodID); err != nil {
			return nil, errors.New("head of department not found")
		}
		dept.HeadOfDepartmentID = &hodID
	}

	if err := s.deptRepo.Create(dept); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return s.toResponse(dept), nil
}

// GetByID gets a department by ID
func (s *DepartmentService) GetByID(id, institutionID uuid.UUID) (*response.DepartmentResponse, error) {
	dept, err := s.deptRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(dept)

	// Get staff count
	count, err := s.deptRepo.GetStaffCount(id)
	if err == nil {
		resp.StaffCount = count
	}

	return resp, nil
}

// GetAll gets all departments with filters
func (s *DepartmentService) GetAll(filter repository.DepartmentFilter, params utils.PaginationParams) ([]response.DepartmentResponse, utils.Pagination, error) {
	departments, total, err := s.deptRepo.FindAll(filter, params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.DepartmentResponse
	for _, dept := range departments {
		resp := s.toResponse(&dept)
		// Get staff count for each department
		count, err := s.deptRepo.GetStaffCount(dept.ID)
		if err == nil {
			resp.StaffCount = count
		}
		responses = append(responses, *resp)
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return responses, pagination, nil
}

// Update updates a department
func (s *DepartmentService) Update(id uuid.UUID, req *request.UpdateDepartmentRequest, institutionID uuid.UUID) (*response.DepartmentResponse, error) {
	dept, err := s.deptRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return nil, err
	}

	// Update name if provided
	if req.Name != "" && req.Name != dept.Name {
		exists, err := s.deptRepo.NameExists(req.Name, institutionID, &id)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, errors.New("department with this name already exists")
		}
		dept.Name = req.Name
	}

	if req.Description != "" {
		dept.Description = req.Description
	}

	// Update head of department if provided
	if req.HeadOfDepartmentID != "" {
		hodID, err := uuid.Parse(req.HeadOfDepartmentID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		if _, err := s.teacherRepo.FindByID(hodID); err != nil {
			return nil, errors.New("head of department not found")
		}
		dept.HeadOfDepartmentID = &hodID
	}

	if err := s.deptRepo.Update(dept); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return s.toResponse(dept), nil
}

// Delete deletes a department
func (s *DepartmentService) Delete(id, institutionID uuid.UUID) error {
	// Verify it exists and belongs to the institution
	_, err := s.deptRepo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return err
	}

	// Check if department has staff
	count, err := s.deptRepo.GetStaffCount(id)
	if err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}
	if count > 0 {
		return errors.New("cannot delete department with staff members")
	}

	return s.deptRepo.Delete(id)
}

// GetDepartmentStaff gets all staff in a department
func (s *DepartmentService) GetDepartmentStaff(deptID, institutionID uuid.UUID) ([]response.TeacherBrief, error) {
	// Verify department exists and belongs to the institution
	_, err := s.deptRepo.FindByIDWithInstitution(deptID, institutionID)
	if err != nil {
		return nil, err
	}

	teachers, err := s.deptRepo.GetDepartmentStaff(deptID)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.TeacherBrief
	for _, t := range teachers {
		brief := response.TeacherBrief{
			ID: t.ID,
		}
		if t.User != nil && t.User.Profile != nil {
			brief.FirstName = t.User.Profile.FirstName
			brief.LastName = t.User.Profile.LastName
		}
		responses = append(responses, brief)
	}

	return responses, nil
}

// toResponse converts a model to response
func (s *DepartmentService) toResponse(dept *models.Department) *response.DepartmentResponse {
	resp := &response.DepartmentResponse{
		ID:            dept.ID,
		InstitutionID: dept.InstitutionID,
		Name:          dept.Name,
		Description:   dept.Description,
		CreatedAt:     dept.CreatedAt,
		UpdatedAt:     dept.UpdatedAt,
	}

	if dept.HeadOfDepartmentID != nil {
		resp.HeadOfDepartmentID = dept.HeadOfDepartmentID
		if dept.HeadOfDepartment != nil {
			resp.HeadOfDepartment = &response.TeacherBrief{
				ID: dept.HeadOfDepartment.ID,
			}
			if dept.HeadOfDepartment.User != nil && dept.HeadOfDepartment.User.Profile != nil {
				resp.HeadOfDepartment.FirstName = dept.HeadOfDepartment.User.Profile.FirstName
				resp.HeadOfDepartment.LastName = dept.HeadOfDepartment.User.Profile.LastName
			}
		}
	}

	return resp
}
