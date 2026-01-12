package service

import (
	"errors"

	"campus-core/internal/dto/request"
	"campus-core/internal/dto/response"
	"campus-core/internal/models"
	"campus-core/internal/repository"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ParentService handles parent management logic
type ParentService struct {
	repo       *repository.ParentRepository
	userRepo   *repository.UserRepository
	db         *gorm.DB
	jwtManager *utils.JWTManager
}

func NewParentService(repo *repository.ParentRepository, userRepo *repository.UserRepository, db *gorm.DB, jwtManager *utils.JWTManager) *ParentService {
	return &ParentService{
		repo:       repo,
		userRepo:   userRepo,
		db:         db,
		jwtManager: jwtManager,
	}
}

// CreateParent creates a new parent
func (s *ParentService) CreateParent(req *request.CreateParentRequest, creatorInstitutionID string) (*response.UserResponse, error) {
	if req.InstitutionID == "" {
		req.InstitutionID = creatorInstitutionID
	}
	if req.InstitutionID == "" {
		return nil, errors.New("institution_id is required")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	institutionID, _ := uuid.Parse(req.InstitutionID)

	var parentUser *models.User
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create User
		user := &models.User{
			BaseModel:    models.BaseModel{ID: uuid.New()},
			Email:        req.Email,
			Phone:        req.Phone,
			PasswordHash: hashedPassword,
			Role:         models.RoleParent,
			IsActive:     true,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 2. Create Profile
		profile := &models.UserProfile{
			BaseModel:     models.BaseModel{ID: uuid.New()},
			UserID:        user.ID,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			InstitutionID: &institutionID,
			Occupation:    req.Occupation,
		}
		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		user.Profile = profile
		parentUser = user

		// 3. Create Parent
		parent := &models.Parent{
			TenantBaseModel: models.TenantBaseModel{
				BaseModel:     models.BaseModel{ID: uuid.New()},
				InstitutionID: institutionID,
			},
			UserID:           user.ID,
			Occupation:       req.Occupation,
			OfficeAddress:    req.OfficeAddress,
			EmergencyContact: req.EmergencyContact,
		}
		if err := tx.Create(parent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	resp := response.UserResponse{
		ID:       parentUser.ID,
		Email:    parentUser.Email,
		Phone:    parentUser.Phone,
		Role:     parentUser.Role,
		IsActive: parentUser.IsActive,
		Profile: &response.ProfileResponse{
			ID:            parentUser.Profile.ID,
			FirstName:     parentUser.Profile.FirstName,
			LastName:      parentUser.Profile.LastName,
			InstitutionID: parentUser.Profile.InstitutionID,
		},
	}

	return &resp, nil
}

// GetAllParents returns all parents
func (s *ParentService) GetAllParents(institutionID string, params utils.PaginationParams) ([]response.UserResponse, utils.Pagination, error) {
	parents, total, err := s.repo.FindAll(institutionID, params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.UserResponse
	for _, p := range parents {
		if p.User.ID != uuid.Nil {
			responses = append(responses, response.UserResponse{
				ID:       p.User.ID,
				Email:    p.User.Email,
				Phone:    p.User.Phone,
				Role:     p.User.Role,
				IsActive: p.User.IsActive,
				Profile: &response.ProfileResponse{
					ID:            p.User.Profile.ID,
					FirstName:     p.User.Profile.FirstName,
					LastName:      p.User.Profile.LastName,
					InstitutionID: p.User.Profile.InstitutionID,
				},
			})
		}
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return responses, pagination, nil
}

// GetParent gets a parent by ID
func (s *ParentService) GetParent(id uuid.UUID) (*response.UserResponse, error) {
	parent, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := response.UserResponse{
		ID:       parent.User.ID,
		Email:    parent.User.Email,
		Phone:    parent.User.Phone,
		Role:     parent.User.Role,
		IsActive: parent.User.IsActive,
		Profile: &response.ProfileResponse{
			ID:            parent.User.Profile.ID,
			FirstName:     parent.User.Profile.FirstName,
			LastName:      parent.User.Profile.LastName,
			InstitutionID: parent.User.Profile.InstitutionID,
		},
	}
	return &resp, nil
}
