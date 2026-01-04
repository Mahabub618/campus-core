package response

import (
	"time"

	"github.com/google/uuid"
)

// LoginResponse represents a successful login response
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresAt    time.Time    `json:"expires_at"`
	User         UserResponse `json:"user"`
}

// TokenResponse represents a token refresh response
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID          uuid.UUID        `json:"id"`
	Email       string           `json:"email,omitempty"`
	Phone       string           `json:"phone,omitempty"`
	Role        string           `json:"role"`
	IsActive    bool             `json:"is_active"`
	LastLoginAt *time.Time       `json:"last_login_at,omitempty"`
	Profile     *ProfileResponse `json:"profile,omitempty"`
}

// ProfileResponse represents user profile data in responses
type ProfileResponse struct {
	ID              uuid.UUID  `json:"id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	FullName        string     `json:"full_name"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"`
	Gender          string     `json:"gender,omitempty"`
	Address         string     `json:"address,omitempty"`
	ProfileImageURL string     `json:"profile_image_url,omitempty"`
	InstitutionID   *uuid.UUID `json:"institution_id,omitempty"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}
