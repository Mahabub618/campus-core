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

// AcademicYearService handles academic year business logic
type AcademicYearService struct {
	repo *repository.AcademicYearRepository
}

// NewAcademicYearService creates a new academic year service
func NewAcademicYearService(repo *repository.AcademicYearRepository) *AcademicYearService {
	return &AcademicYearService{repo: repo}
}

// Create creates a new academic year
func (s *AcademicYearService) Create(req *request.CreateAcademicYearRequest, institutionID uuid.UUID) (*response.AcademicYearResponse, error) {
	// Validate date range
	if req.EndDate.Before(req.StartDate) || req.EndDate.Equal(req.StartDate) {
		return nil, errors.New("end date must be after start date")
	}

	// Check if name already exists
	exists, err := s.repo.NameExists(req.Name, institutionID, nil)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}
	if exists {
		return nil, errors.New("academic year with this name already exists")
	}

	ay := &models.AcademicYear{
		InstitutionID: institutionID,
		Name:          req.Name,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		IsCurrent:     req.IsCurrent,
		Description:   req.Description,
	}

	if err := s.repo.Create(ay); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// If this is set as current, update other academic years
	if req.IsCurrent {
		if err := s.repo.SetCurrent(ay.ID, institutionID); err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
	}

	return s.toResponse(ay), nil
}

// GetByID gets an academic year by ID
func (s *AcademicYearService) GetByID(id, institutionID uuid.UUID) (*response.AcademicYearResponse, error) {
	ay, err := s.repo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return nil, err
	}
	return s.toResponse(ay), nil
}

// GetAll gets all academic years with filters
func (s *AcademicYearService) GetAll(filter repository.AcademicYearFilter, params utils.PaginationParams) ([]response.AcademicYearResponse, utils.Pagination, error) {
	academicYears, total, err := s.repo.FindAll(filter, params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.AcademicYearResponse
	for _, ay := range academicYears {
		responses = append(responses, *s.toResponse(&ay))
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return responses, pagination, nil
}

// GetCurrent gets the current academic year for an institution
func (s *AcademicYearService) GetCurrent(institutionID uuid.UUID) (*response.AcademicYearResponse, error) {
	ay, err := s.repo.FindCurrent(institutionID)
	if err != nil {
		return nil, err
	}
	return s.toResponse(ay), nil
}

// Update updates an academic year
func (s *AcademicYearService) Update(id uuid.UUID, req *request.UpdateAcademicYearRequest, institutionID uuid.UUID) (*response.AcademicYearResponse, error) {
	ay, err := s.repo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" && req.Name != ay.Name {
		exists, err := s.repo.NameExists(req.Name, institutionID, &id)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, errors.New("academic year with this name already exists")
		}
		ay.Name = req.Name
	}

	if req.StartDate != nil {
		ay.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		ay.EndDate = *req.EndDate
	}
	if req.Description != "" {
		ay.Description = req.Description
	}

	// Validate date range
	if ay.EndDate.Before(ay.StartDate) || ay.EndDate.Equal(ay.StartDate) {
		return nil, errors.New("end date must be after start date")
	}

	if err := s.repo.Update(ay); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// Handle is_current update separately
	if req.IsCurrent != nil && *req.IsCurrent {
		if err := s.repo.SetCurrent(id, institutionID); err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		ay.IsCurrent = true
	}

	return s.toResponse(ay), nil
}

// Delete deletes an academic year
func (s *AcademicYearService) Delete(id, institutionID uuid.UUID) error {
	// Verify it exists and belongs to the institution
	_, err := s.repo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// Activate sets an academic year as current
func (s *AcademicYearService) Activate(id, institutionID uuid.UUID) error {
	// Verify it exists and belongs to the institution
	_, err := s.repo.FindByIDWithInstitution(id, institutionID)
	if err != nil {
		return err
	}

	return s.repo.SetCurrent(id, institutionID)
}

// toResponse converts a model to response
func (s *AcademicYearService) toResponse(ay *models.AcademicYear) *response.AcademicYearResponse {
	return &response.AcademicYearResponse{
		ID:            ay.ID,
		InstitutionID: ay.InstitutionID,
		Name:          ay.Name,
		StartDate:     ay.StartDate,
		EndDate:       ay.EndDate,
		IsCurrent:     ay.IsCurrent,
		Description:   ay.Description,
		CreatedAt:     ay.CreatedAt,
		UpdatedAt:     ay.UpdatedAt,
	}
}
