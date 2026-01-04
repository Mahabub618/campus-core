package database

import (
	"campus-core/internal/models"
	"campus-core/internal/utils"
	"campus-core/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Seeder handles database seeding
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// SeedAll runs all seed functions
func (s *Seeder) SeedAll() error {
	logger.Info("Starting database seeding...")

	if err := s.SeedSuperAdmin(); err != nil {
		return err
	}

	logger.Info("Database seeding completed successfully")
	return nil
}

// SeedSuperAdmin creates a default super admin user if not exists
func (s *Seeder) SeedSuperAdmin() error {
	var count int64
	s.db.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&count)

	if count > 0 {
		logger.Info("Super admin already exists, skipping seed")
		return nil
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword("Admin@123")
	if err != nil {
		return err
	}

	// Create super admin user
	superAdmin := &models.User{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		Email:        "superadmin@campus.local",
		PasswordHash: hashedPassword,
		Role:         models.RoleSuperAdmin,
		IsActive:     true,
	}

	if err := s.db.Create(superAdmin).Error; err != nil {
		return err
	}

	// Create profile for super admin
	profile := &models.UserProfile{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		UserID:    superAdmin.ID,
		FirstName: "Super",
		LastName:  "Admin",
	}

	if err := s.db.Create(profile).Error; err != nil {
		return err
	}

	logger.Info("Super admin created successfully",
		zap.String("email", superAdmin.Email),
		zap.String("password", "Admin@123"),
	)

	return nil
}

// SeedTestInstitution creates a test institution if not exists
func (s *Seeder) SeedTestInstitution() error {
	var count int64
	s.db.Model(&models.Institution{}).Where("code = ?", "TEST001").Count(&count)

	if count > 0 {
		logger.Info("Test institution already exists, skipping seed")
		return nil
	}

	institution := &models.Institution{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		Name:            "Test School",
		Code:            "TEST001",
		Address:         "123 Test Street",
		Phone:           "+1234567890",
		Email:           "test@school.com",
		PrincipalName:   "John Doe",
		EstablishedYear: 2000,
		IsActive:        true,
	}

	if err := s.db.Create(institution).Error; err != nil {
		return err
	}

	logger.Info("Test institution created successfully",
		zap.String("name", institution.Name),
		zap.String("code", institution.Code),
	)

	return nil
}
