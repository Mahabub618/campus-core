package handler

import (
	"net/http"

	"campus-core/internal/dto/request"
	"campus-core/internal/middleware"
	"campus-core/internal/service"
	"campus-core/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with email/phone and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "Login credentials"
// @Success 200 {object} utils.APIResponse{data=response.LoginResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err)
		return
	}

	utils.OK(c, "Login successful", resp)
}

// Register handles user registration (admin only)
// @Summary Register new user
// @Description Create a new user account (admin only)
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.RegisterRequest true "User registration data"
// @Success 201 {object} utils.APIResponse{data=response.UserResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.Created(c, "User registered successfully", resp)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} utils.APIResponse{data=response.TokenResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.authService.RefreshToken(&req)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err)
		return
	}

	utils.OK(c, "Token refreshed successfully", resp)
}

// Logout handles user logout
// @Summary User logout
// @Description Invalidate user's refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ErrTokenMissing)
		return
	}

	if err := h.authService.Logout(userID); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.OK(c, "Logged out successfully", nil)
}

// ForgotPassword handles password reset request
// @Summary Forgot password
// @Description Request password reset email
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.ForgotPasswordRequest true "Email address"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.authService.ForgotPassword(&req); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Always return success to prevent email enumeration
	utils.OK(c, "If the email exists, a password reset link has been sent", nil)
}

// ResetPassword handles password reset with token
// @Summary Reset password
// @Description Reset password using reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.authService.ResetPassword(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Password reset successfully", nil)
}

// ChangePassword handles password change for authenticated users
// @Summary Change password
// @Description Change password for authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.ChangePasswordRequest true "Old and new password"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ErrTokenMissing)
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.authService.ChangePassword(userID, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.OK(c, "Password changed successfully", nil)
}

// GetMe returns the current authenticated user
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse{data=response.UserResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ErrTokenMissing)
		return
	}

	resp, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.OK(c, "", resp)
}
