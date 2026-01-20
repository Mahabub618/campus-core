package service

import (
	"errors"
	"time"

	"campus-core/internal/dto/request"
	"campus-core/internal/dto/response"
	"campus-core/internal/models"
	"campus-core/internal/repository"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountantService handles accountant management logic
type AccountantService struct {
	repo       *repository.AccountantRepository
	userRepo   *repository.UserRepository
	db         *gorm.DB
	jwtManager *utils.JWTManager
}

func NewAccountantService(repo *repository.AccountantRepository, userRepo *repository.UserRepository, db *gorm.DB, jwtManager *utils.JWTManager) *AccountantService {
	return &AccountantService{
		repo:       repo,
		userRepo:   userRepo,
		db:         db,
		jwtManager: jwtManager,
	}
}

// CreateAccountant creates a new accountant
func (s *AccountantService) CreateAccountant(req *request.CreateAccountantRequest, creatorInstitutionID string) (*response.UserResponse, error) {
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

	var accountantUser *models.User
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create User
		user := &models.User{
			BaseModel:    models.BaseModel{ID: uuid.New()},
			Email:        req.Email,
			Phone:        req.Phone,
			PasswordHash: hashedPassword,
			Role:         models.RoleAccountant,
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
			Occupation:    "Accountant",
		}
		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		user.Profile = profile
		accountantUser = user

		// 3. Create Accountant
		joiningDate, _ := time.Parse("2006-01-02", req.JoiningDate)
		accountant := &models.Accountant{
			TenantBaseModel: models.TenantBaseModel{
				BaseModel:     models.BaseModel{ID: uuid.New()},
				InstitutionID: institutionID,
			},
			UserID:        user.ID,
			JoiningDate:   &joiningDate,
			Qualification: req.Qualification,
		}
		if err := tx.Create(accountant).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	resp := response.UserResponse{
		ID:       accountantUser.ID,
		Email:    accountantUser.Email,
		Phone:    accountantUser.Phone,
		Role:     accountantUser.Role,
		IsActive: accountantUser.IsActive,
		Profile: &response.ProfileResponse{
			ID:            accountantUser.Profile.ID,
			FirstName:     accountantUser.Profile.FirstName,
			LastName:      accountantUser.Profile.LastName,
			InstitutionID: accountantUser.Profile.InstitutionID,
		},
	}

	return &resp, nil
}

// GetAllAccountants returns all accountants
func (s *AccountantService) GetAllAccountants(institutionID string, params utils.PaginationParams) ([]response.UserResponse, utils.Pagination, error) {
	accountants, total, err := s.repo.FindAll(institutionID, params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.UserResponse
	for _, a := range accountants {
		if a.User.ID != uuid.Nil {
			responses = append(responses, response.UserResponse{
				ID:       a.User.ID,
				Email:    a.User.Email,
				Phone:    a.User.Phone,
				Role:     a.User.Role,
				IsActive: a.User.IsActive,
				Profile: &response.ProfileResponse{
					ID:            a.User.Profile.ID,
					FirstName:     a.User.Profile.FirstName,
					LastName:      a.User.Profile.LastName,
					InstitutionID: a.User.Profile.InstitutionID,
				},
			})
		}
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return responses, pagination, nil
}

// GetAccountant gets an accountant by ID
func (s *AccountantService) GetAccountant(id uuid.UUID) (*response.UserResponse, error) {
	accountant, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := response.UserResponse{
		ID:       accountant.User.ID,
		Email:    accountant.User.Email,
		Phone:    accountant.User.Phone,
		Role:     accountant.User.Role,
		IsActive: accountant.User.IsActive,
		Profile: &response.ProfileResponse{
			ID:            accountant.User.Profile.ID,
			FirstName:     accountant.User.Profile.FirstName,
			LastName:      accountant.User.Profile.LastName,
			InstitutionID: accountant.User.Profile.InstitutionID,
		},
	}
	return &resp, nil
}

// UpdateAccountant updates an accountant
func (s *AccountantService) UpdateAccountant(id uuid.UUID, req *request.UpdateAccountantRequest, institutionID string) (*response.UserResponse, error) {
	accountant, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Verify tenant access
	if institutionID != "" && accountant.InstitutionID.String() != institutionID {
		return nil, utils.ErrCrossTenantAccess
	}

	// Update user fields
	if req.Email != "" && req.Email != accountant.User.Email {
		var count int64
		if err := s.db.Model(&models.User{}).Where("email = ? AND id != ?", req.Email, accountant.User.ID).Count(&count).Error; err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if count > 0 {
			return nil, utils.ErrEmailAlreadyExists
		}
		accountant.User.Email = req.Email
	}

	if req.Phone != "" {
		accountant.User.Phone = req.Phone
	}

	if req.IsActive != nil {
		accountant.User.IsActive = *req.IsActive
	}

	// Update profile fields
	if accountant.User.Profile != nil {
		if req.FirstName != "" {
			accountant.User.Profile.FirstName = req.FirstName
		}
		if req.LastName != "" {
			accountant.User.Profile.LastName = req.LastName
		}
	}

	// Update accountant-specific fields
	if req.Qualification != "" {
		accountant.Qualification = req.Qualification
	}

	// Save changes in transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(accountant.User).Error; err != nil {
			return err
		}
		if accountant.User.Profile != nil {
			if err := tx.Save(accountant.User.Profile).Error; err != nil {
				return err
			}
		}
		if err := tx.Save(accountant).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	resp := response.UserResponse{
		ID:       accountant.User.ID,
		Email:    accountant.User.Email,
		Phone:    accountant.User.Phone,
		Role:     accountant.User.Role,
		IsActive: accountant.User.IsActive,
		Profile: &response.ProfileResponse{
			ID:            accountant.User.Profile.ID,
			FirstName:     accountant.User.Profile.FirstName,
			LastName:      accountant.User.Profile.LastName,
			InstitutionID: accountant.User.Profile.InstitutionID,
		},
	}
	return &resp, nil
}
