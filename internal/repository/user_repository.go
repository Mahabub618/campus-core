package repository

import (
	"errors"
	"time"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Profile").First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Profile").First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByPhone finds a user by phone
func (r *UserRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Profile").First(&user, "phone = ?", phone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmailOrPhone finds a user by email or phone
func (r *UserRepository) FindByEmailOrPhone(identifier string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Profile").First(&user, "email = ? OR phone = ?", identifier, identifier).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// UpdateLastLogin updates the last login time
func (r *UserRepository) UpdateLastLogin(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}

// SaveRefreshToken saves or updates the refresh token for a user
func (r *UserRepository) SaveRefreshToken(id uuid.UUID, token string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("refresh_token", token).Error
}

// InvalidateRefreshToken clears the refresh token for a user
func (r *UserRepository) InvalidateRefreshToken(id uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("refresh_token", "").Error
}

// FindByRefreshToken finds a user by refresh token
func (r *UserRepository) FindByRefreshToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Profile").First(&user, "refresh_token = ?", token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRefreshTokenInvalid
		}
		return nil, err
	}
	return &user, nil
}

// SaveResetToken saves a password reset token
func (r *UserRepository) SaveResetToken(id uuid.UUID, token string, expiry time.Time) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"reset_token":        token,
		"reset_token_expiry": expiry,
	}).Error
}

// FindByResetToken finds a user by reset token
func (r *UserRepository) FindByResetToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "reset_token = ? AND reset_token_expiry > ?", token, time.Now()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrResetTokenInvalid
		}
		return nil, err
	}
	return &user, nil
}

// ClearResetToken clears the reset token after use
func (r *UserRepository) ClearResetToken(id uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"reset_token":        "",
		"reset_token_expiry": nil,
	}).Error
}

// UpdatePassword updates the user's password
func (r *UserRepository) UpdatePassword(id uuid.UUID, passwordHash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password_hash", passwordHash).Error
}

// EmailExists checks if an email is already registered
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// PhoneExists checks if a phone is already registered
func (r *UserRepository) PhoneExists(phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error
	return count > 0, err
}

// CreateWithProfile creates a user with profile in a transaction
func (r *UserRepository) CreateWithProfile(user *models.User, profile *models.UserProfile) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		profile.UserID = user.ID
		if err := tx.Create(profile).Error; err != nil {
			return err
		}

		return nil
	})
}
