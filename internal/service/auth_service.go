package service

import (
	"campus-core/internal/dto/request"
	"campus-core/internal/dto/response"
	"campus-core/internal/middleware"
	"campus-core/internal/models"
	"campus-core/internal/repository"
	"campus-core/internal/utils"
	"campus-core/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *utils.JWTManager
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, jwtManager *utils.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	var user *models.User
	var err error

	// Find user by email or phone
	if req.Email != "" {
		user, err = s.userRepo.FindByEmail(req.Email)
	} else if req.Phone != "" {
		user, err = s.userRepo.FindByPhone(req.Phone)
	} else {
		return nil, utils.ErrInvalidCredentials
	}

	if err != nil {
		logger.Debug("User not found during login", zap.String("email", req.Email))
		return nil, utils.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, utils.ErrAccountDisabled
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, utils.ErrInvalidCredentials
	}

	// Get institution ID from profile if available
	institutionID := ""
	if user.Profile != nil && user.Profile.InstitutionID != nil {
		institutionID = user.Profile.InstitutionID.String()
	}

	// Get permissions for the user's role
	permissions := middleware.GetPermissionsForRole(user.Role)

	// Generate access token
	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
		institutionID,
		permissions,
	)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// Generate refresh token
	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// Save refresh token to database
	if err := s.userRepo.SaveRefreshToken(user.ID, refreshToken); err != nil {
		logger.Error("Failed to save refresh token", zap.Error(err))
	}

	// Update last login time
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		logger.Error("Failed to update last login", zap.Error(err))
	}

	return &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
		User:         s.toUserResponse(user),
	}, nil
}

// Register creates a new user (admin only)
func (s *AuthService) Register(req *request.RegisterRequest) (*response.UserResponse, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}
	if exists {
		return nil, utils.ErrEmailAlreadyExists
	}

	// Check if phone already exists (if provided)
	if req.Phone != "" {
		exists, err = s.userRepo.PhoneExists(req.Phone)
		if err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if exists {
			return nil, utils.ErrPhoneAlreadyExists
		}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// Create user
	user := &models.User{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		IsActive:     true,
	}

	// Create profile
	profile := &models.UserProfile{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Set institution ID if provided
	if req.InstitutionID != "" {
		instID, err := uuid.Parse(req.InstitutionID)
		if err != nil {
			return nil, utils.ErrInvalidUUID
		}
		profile.InstitutionID = &instID
	}

	// Create user with profile
	if err := s.userRepo.CreateWithProfile(user, profile); err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	user.Profile = profile
	resp := s.toUserResponse(user)
	return &resp, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *AuthService) RefreshToken(req *request.RefreshTokenRequest) (*response.TokenResponse, error) {
	// Validate refresh token
	userID, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Find user and verify refresh token matches
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	if user.RefreshToken != req.RefreshToken {
		return nil, utils.ErrRefreshTokenInvalid
	}

	if !user.IsActive {
		return nil, utils.ErrAccountDisabled
	}

	// Get institution ID
	institutionID := ""
	if user.Profile != nil && user.Profile.InstitutionID != nil {
		institutionID = user.Profile.InstitutionID.String()
	}

	// Get permissions
	permissions := middleware.GetPermissionsForRole(user.Role)

	// Generate new access token
	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
		institutionID,
		permissions,
	)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// Generate new refresh token
	refreshToken, _, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// Save new refresh token
	if err := s.userRepo.SaveRefreshToken(user.ID, refreshToken); err != nil {
		logger.Error("Failed to save refresh token", zap.Error(err))
	}

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout invalidates the user's refresh token
func (s *AuthService) Logout(userID uuid.UUID) error {
	return s.userRepo.InvalidateRefreshToken(userID)
}

// ForgotPassword initiates the password reset process
func (s *AuthService) ForgotPassword(req *request.ForgotPasswordRequest) error {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists
		logger.Debug("Forgot password for non-existent email", zap.String("email", req.Email))
		return nil
	}

	// Generate reset token
	resetToken, expiry, err := s.jwtManager.GenerateResetToken(user.ID, user.Email)
	if err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}

	// Save reset token
	if err := s.userRepo.SaveResetToken(user.ID, resetToken, expiry); err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}

	// TODO: Send email with reset token
	// For now, just log it (in development)
	logger.Info("Password reset token generated",
		zap.String("email", user.Email),
		zap.String("token", resetToken),
		zap.Time("expiry", expiry),
	)

	return nil
}

// ResetPassword resets the user's password using a reset token
func (s *AuthService) ResetPassword(req *request.ResetPasswordRequest) error {
	// Validate reset token
	userID, err := s.jwtManager.ValidateResetToken(req.Token)
	if err != nil {
		return err
	}

	// Find user and verify token matches
	user, err := s.userRepo.FindByResetToken(req.Token)
	if err != nil {
		return err
	}

	if user.ID != userID {
		return utils.ErrResetTokenInvalid
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(user.ID, hashedPassword); err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}

	// Clear reset token
	if err := s.userRepo.ClearResetToken(user.ID); err != nil {
		logger.Error("Failed to clear reset token", zap.Error(err))
	}

	// Invalidate all refresh tokens
	if err := s.userRepo.InvalidateRefreshToken(user.ID); err != nil {
		logger.Error("Failed to invalidate refresh token", zap.Error(err))
	}

	return nil
}

// ChangePassword changes the user's password
func (s *AuthService) ChangePassword(userID uuid.UUID, req *request.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !utils.CheckPassword(req.OldPassword, user.PasswordHash) {
		return utils.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}

	return nil
}

// GetCurrentUser returns the current user's information
func (s *AuthService) GetCurrentUser(userID uuid.UUID) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	resp := s.toUserResponse(user)
	return &resp, nil
}

// toUserResponse converts a user model to response DTO
func (s *AuthService) toUserResponse(user *models.User) response.UserResponse {
	resp := response.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Phone:       user.Phone,
		Role:        user.Role,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
	}

	if user.Profile != nil {
		resp.Profile = &response.ProfileResponse{
			ID:              user.Profile.ID,
			FirstName:       user.Profile.FirstName,
			LastName:        user.Profile.LastName,
			FullName:        user.Profile.FullName(),
			DateOfBirth:     user.Profile.DateOfBirth,
			Gender:          user.Profile.Gender,
			Address:         user.Profile.Address,
			ProfileImageURL: user.Profile.ProfileImageURL,
			InstitutionID:   user.Profile.InstitutionID,
		}
	}

	return resp
}
