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

// UpdateUser updates a user (Admin function)
func (s *UserService) UpdateUser(id uuid.UUID, req *request.UpdateUserRequest, creatorRole string, creatorInstitutionID string) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Security: Verify tenant access for non-super admins
	if creatorRole != models.RoleSuperAdmin {
		if user.Profile != nil && user.Profile.InstitutionID != nil {
			if user.Profile.InstitutionID.String() != creatorInstitutionID {
				return nil, utils.ErrCrossTenantAccess
			}
		}
		// Admin cannot update Super Admins
		if user.Role == models.RoleSuperAdmin {
			return nil, utils.ErrActionNotPermitted
		}
	}

	// Update email if provided and changed
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.repo.EmailExists(req.Email)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, utils.ErrEmailAlreadyExists
		}
		user.Email = req.Email
	}

	// Update phone if provided and changed
	if req.Phone != "" && req.Phone != user.Phone {
		exists, err := s.repo.PhoneExists(req.Phone)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, utils.ErrPhoneAlreadyExists
		}
		user.Phone = req.Phone
	}

	// Update active status if provided
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update profile fields
	if user.Profile != nil {
		if req.FirstName != "" {
			user.Profile.FirstName = req.FirstName
		}
		if req.LastName != "" {
			user.Profile.LastName = req.LastName
		}
	}

	if err := s.repo.Update(user); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	resp := s.authService.toUserResponse(user)
	return &resp, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uuid.UUID, creatorRole string, creatorInstitutionID string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Security: Verify tenant access for non-super admins
	if creatorRole != models.RoleSuperAdmin {
		if user.Profile != nil && user.Profile.InstitutionID != nil {
			if user.Profile.InstitutionID.String() != creatorInstitutionID {
				return utils.ErrCrossTenantAccess
			}
		}
		// Admin cannot delete Super Admins
		if user.Role == models.RoleSuperAdmin {
			return utils.ErrActionNotPermitted
		}
	}

	return s.repo.Delete(id)
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

// UpdateAvatar updates the user's avatar
func (s *UserService) UpdateAvatar(userID uuid.UUID, avatarURL string) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Profile == nil {
		user.Profile = &models.UserProfile{UserID: userID}
	}

	user.Profile.ProfileImageURL = avatarURL

	if err := s.repo.Update(user); err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	resp := s.authService.toUserResponse(user)
	return &resp, nil
}

// UpdatePassword updates the user's password
func (s *UserService) UpdatePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !utils.CheckPassword(oldPassword, user.PasswordHash) {
		return utils.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}

	return s.repo.UpdatePassword(userID, hashedPassword)
}
