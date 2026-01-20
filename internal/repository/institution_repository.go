package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InstitutionRepository handles database operations for institutions
type InstitutionRepository struct {
	db *gorm.DB
}

// NewInstitutionRepository creates a new institution repository
func NewInstitutionRepository(db *gorm.DB) *InstitutionRepository {
	return &InstitutionRepository{db: db}
}

// Create creates a new institution
func (r *InstitutionRepository) Create(institution *models.Institution) error {
	return r.db.Create(institution).Error
}

// FindByID finds an institution by ID
func (r *InstitutionRepository) FindByID(id uuid.UUID) (*models.Institution, error) {
	var institution models.Institution
	if err := r.db.First(&institution, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrInstitutionNotFound
		}
		return nil, err
	}
	return &institution, nil
}

// FindByCode finds an institution by code
func (r *InstitutionRepository) FindByCode(code string) (*models.Institution, error) {
	var institution models.Institution
	if err := r.db.First(&institution, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrInstitutionNotFound
		}
		return nil, err
	}
	return &institution, nil
}

// Update updates an institution
func (r *InstitutionRepository) Update(institution *models.Institution) error {
	return r.db.Save(institution).Error
}

// Delete deletes an institution
func (r *InstitutionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Institution{}, "id = ?", id).Error
}

// FindAll returns a list of institutions with pagination
func (r *InstitutionRepository) FindAll(params utils.PaginationParams) ([]models.Institution, int64, error) {
	var institutions []models.Institution
	var total int64

	// Count total
	if err := r.db.Model(&models.Institution{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err := r.db.Scopes(utils.Paginate(params)).Find(&institutions).Error
	if err != nil {
		return nil, 0, err
	}

	return institutions, total, nil
}

// GetStats returns statistics for an institution
func (r *InstitutionRepository) GetStats(id uuid.UUID) (*models.InstitutionStats, error) {
	var stats models.InstitutionStats
	stats.InstitutionID = id

	// Count students
	if err := r.db.Model(&models.Student{}).Where("institution_id = ?", id).Count(&stats.TotalStudents).Error; err != nil {
		return nil, err
	}

	// Count teachers
	if err := r.db.Model(&models.Teacher{}).Where("institution_id = ?", id).Count(&stats.TotalTeachers).Error; err != nil {
		return nil, err
	}

	// Count parents
	if err := r.db.Model(&models.Parent{}).Where("institution_id = ?", id).Count(&stats.TotalParents).Error; err != nil {
		return nil, err
	}

	// Count users
	// Join with user_profiles to filter by institution
	// Note: Users table doesn't have institution_id directly, logic might need adjustment if users belong to multiple
	// But in our schema, UserProfile has InstitutionID
	if err := r.db.Model(&models.UserProfile{}).Where("institution_id = ?", id).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// CodeExists checks if a code already exists
func (r *InstitutionRepository) CodeExists(code string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Institution{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

// GetAdmins returns all admin users for an institution
func (r *InstitutionRepository) GetAdmins(institutionID uuid.UUID) ([]models.User, error) {
	var users []models.User

	err := r.db.Preload("Profile").
		Joins("INNER JOIN user_profiles ON user_profiles.user_id = users.id").
		Where("user_profiles.institution_id = ? AND users.role = ?", institutionID, models.RoleAdmin).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// CreateAdmin creates a new admin user for an institution
func (r *InstitutionRepository) CreateAdmin(institutionID uuid.UUID, email, firstName, lastName, password, phone string) (*models.User, error) {
	// Check if email already exists
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, utils.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	var user *models.User
	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Create user
		user = &models.User{
			BaseModel: models.BaseModel{
				ID: uuid.New(),
			},
			Email:        email,
			Phone:        phone,
			PasswordHash: hashedPassword,
			Role:         models.RoleAdmin,
			IsActive:     true,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// Create profile
		profile := &models.UserProfile{
			BaseModel: models.BaseModel{
				ID: uuid.New(),
			},
			UserID:        user.ID,
			InstitutionID: &institutionID,
			FirstName:     firstName,
			LastName:      lastName,
		}
		if err := tx.Create(profile).Error; err != nil {
			return err
		}

		user.Profile = profile
		return nil
	})

	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	return user, nil
}
