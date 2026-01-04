package request

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required_without=Phone,omitempty,email"`
	Phone    string `json:"phone" binding:"required_without=Email,omitempty"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterRequest represents a user registration request (admin only)
type RegisterRequest struct {
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone" binding:"omitempty,phone"`
	Password      string `json:"password" binding:"required,min=8"`
	Role          string `json:"role" binding:"required,role"`
	FirstName     string `json:"first_name" binding:"required,min=1,max=100"`
	LastName      string `json:"last_name" binding:"required,min=1,max=100"`
	InstitutionID string `json:"institution_id" binding:"omitempty,uuid"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// LogoutRequest represents a logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"omitempty"`
}
