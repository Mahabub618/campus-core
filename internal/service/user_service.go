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

// UserService handles user management business logic
type UserService struct {
	repo        *repository.UserRepository
	instRepo    *repository.InstitutionRepository
	authService *AuthService // Reuse for registration logic including hashing
}

// NewUserService creates a new user service
func NewUserService(repo *repository.UserRepository, instRepo *repository.InstitutionRepository, authService *AuthService) *UserService {
	return &UserService{
		repo:        repo,
		instRepo:    instRepo,
		authService: authService,
	}
}

// CreateUser creates a new user (Admin function)
func (s *UserService) CreateUser(req *request.RegisterRequest, creatorRole string, creatorInstitutionID string) (*response.UserResponse, error) {
	// 1. Validation Logic
	// Super Admin can create any user for any institution (if InstitutionID provided)
	// Admin can only create users for their own institution

	targetInstitutionID := req.InstitutionID
	if targetInstitutionID == "" && creatorRole != models.RoleSuperAdmin {
		targetInstitutionID = creatorInstitutionID
	}

	if creatorRole != models.RoleSuperAdmin {
		if targetInstitutionID != creatorInstitutionID {
			return nil, utils.ErrActionNotPermitted // Cannot create for other institution
		}
		// Admin cannot create Super Admin or other Admins (maybe?)
		// Let's allow Admin to create other Admins for now, or restriction logic here
		if req.Role == models.RoleSuperAdmin {
			return nil, utils.ErrActionNotPermitted
		}
	}

	// Verify Institution exists
	if targetInstitutionID != "" {
		id, err := uuid.Parse(targetInstitutionID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		if _, err := s.instRepo.FindByID(id); err != nil {
			return nil, utils.ErrInstitutionNotFound
		}
	} else if req.Role != models.RoleSuperAdmin {
		// Non-SuperAdmin users must belong to an institution
		// Unless it's a platform-level user? Assumed requirement: All users belong to institution except SA
		return nil, errors.New("institution_id is required")
	}

	// Delegate to AuthService's Register logic (password hashing, etc.)
	// But AuthService.Register takes RegisterRequest defined in its package.
	// Since we imported it, we can pass it.
	// Note: We might want slightly different logic or reuse.
	// Reuse is better to avoid duplication.

	return s.authService.Register(req)
}

// GetUser gets a user by ID
func (s *UserService) GetUser(id uuid.UUID) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := s.authService.toUserResponse(user) // Use helper from auth service or duplicate it
	return &resp, nil
}

// GetAllUsers lists users with filters
func (s *UserService) GetAllUsers(filter repository.UserFilter, params utils.PaginationParams) ([]response.UserResponse, utils.Pagination, error) {
	users, total, err := s.repo.FindAll(filter, params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	var userResponses []response.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, s.authService.toUserResponse(&user))
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return userResponses, pagination, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id uuid.UUID, req map[string]interface{}) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Logic to update user fields...
	// For MVP, focus on Profile updates mainly
	// Updates to User table (email/phone) require uniqueness checks

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	resp := s.authService.toUserResponse(user)
	return &resp, nil
}

// ToggleStatus changes user active status
func (s *UserService) ToggleStatus(id uuid.UUID, isActive bool) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}
	return s.repo.UpdateStatus(id, isActive)
}

// UpdateProfile updates the user's profile
func (s *UserService) UpdateProfile(userID uuid.UUID, firstName, lastName string) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Profile == nil {
		user.Profile = &models.UserProfile{UserID: userID} // Should exist, but safety check
	}

	user.Profile.FirstName = firstName
	user.Profile.LastName = lastName

	if err := s.repo.Update(user); err != nil { // Updates relations too if configured
		return nil, err
	}

	resp := s.authService.toUserResponse(user)
	return &resp, nil
}
