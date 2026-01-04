package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID        uuid.UUID `json:"user_id"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	InstitutionID string    `json:"institution_id,omitempty"`
	Permissions   []string  `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

// TokenType represents the type of token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// JWTManager handles JWT operations
type JWTManager struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secret string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secret:        []byte(secret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateAccessToken generates a new access token
func (m *JWTManager) GenerateAccessToken(userID uuid.UUID, email, role, institutionID string, permissions []string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.accessExpiry)

	claims := &Claims{
		UserID:        userID,
		Email:         email,
		Role:          role,
		InstitutionID: institutionID,
		Permissions:   permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
			Issuer:    "campus-core",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken generates a new refresh token
func (m *JWTManager) GenerateRefreshToken(userID uuid.UUID) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.refreshExpiry)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID.String(),
		Issuer:    "campus-core",
		ID:        uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateAccessToken validates and parses an access token
func (m *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ValidateRefreshToken validates and parses a refresh token
func (m *JWTManager) ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return uuid.Nil, ErrRefreshTokenExpired
		}
		return uuid.Nil, ErrRefreshTokenInvalid
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return uuid.Nil, ErrRefreshTokenInvalid
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, ErrRefreshTokenInvalid
	}

	return userID, nil
}

// GenerateResetToken generates a password reset token
func (m *JWTManager) GenerateResetToken(userID uuid.UUID, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(1 * time.Hour) // Reset token valid for 1 hour

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID.String(),
		Issuer:    "campus-core-reset",
		ID:        uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateResetToken validates a password reset token
func (m *JWTManager) ValidateResetToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return uuid.Nil, ErrResetTokenExpired
		}
		return uuid.Nil, ErrResetTokenInvalid
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid || claims.Issuer != "campus-core-reset" {
		return uuid.Nil, ErrResetTokenInvalid
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, ErrResetTokenInvalid
	}

	return userID, nil
}
