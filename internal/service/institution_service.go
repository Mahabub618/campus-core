package service

import (
	"campus-core/internal/models"
	"campus-core/internal/repository"
	"campus-core/internal/utils"

	"github.com/google/uuid"
)

// InstitutionService handles business logic for institutions
type InstitutionService struct {
	repo *repository.InstitutionRepository
}

// NewInstitutionService creates a new institution service
func NewInstitutionService(repo *repository.InstitutionRepository) *InstitutionService {
	return &InstitutionService{repo: repo}
}

// CreateInstitution creates a new institution
func (s *InstitutionService) Create(institution *models.Institution) error {
	// Check if code exists
	exists, err := s.repo.CodeExists(institution.Code)
	if err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}
	if exists {
		return utils.ErrInstitutionCodeExists
	}

	// Set default ID if not provided (GORM does this, but good to be explicit for logic)
	if institution.ID == uuid.Nil {
		institution.ID = uuid.New()
	}

	if err := s.repo.Create(institution); err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}

	return nil
}

// GetInstitution gets an institution by ID
func (s *InstitutionService) GetByID(id uuid.UUID) (*models.Institution, error) {
	return s.repo.FindByID(id)
}

// UpdateInstitution updates an institution
func (s *InstitutionService) Update(id uuid.UUID, updates map[string]interface{}) (*models.Institution, error) {
	institution, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Prevent code update if it exists
	if code, ok := updates["code"].(string); ok && code != institution.Code {
		exists, err := s.repo.CodeExists(code)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, utils.ErrInstitutionCodeExists
		}
		institution.Code = code
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		institution.Name = name
	}
	if address, ok := updates["address"].(string); ok {
		institution.Address = address
	}
	if phone, ok := updates["phone"].(string); ok {
		institution.Phone = phone
	}
	if email, ok := updates["email"].(string); ok {
		institution.Email = email
	}
	if princ, ok := updates["principal_name"].(string); ok {
		institution.PrincipalName = princ
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		institution.IsActive = isActive
	}

	if err := s.repo.Update(institution); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return institution, nil
}

// DeleteInstitution deletes an institution
func (s *InstitutionService) Delete(id uuid.UUID) error {
	// Check if exists
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}

	// TODO: Check if it has active users/data before deleting?
	// For now, allow soft delete

	if err := s.repo.Delete(id); err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}

	return nil
}

// GetAllInstitutions returns all institutions
func (s *InstitutionService) GetAll(params utils.PaginationParams) ([]models.Institution, utils.Pagination, error) {
	data, total, err := s.repo.FindAll(params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return data, pagination, nil
}

// GetStats returns stats for an institution
func (s *InstitutionService) GetStats(id uuid.UUID) (*models.InstitutionStats, error) {
	// Verify existence
	if _, err := s.repo.FindByID(id); err != nil {
		return nil, err
	}

	stats, err := s.repo.GetStats(id)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return stats, nil
}
