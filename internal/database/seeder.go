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

	// 1. Institutions
	if err := s.SeedInstitutions(); err != nil {
		return err
	}

	// 2. Academic Data
	if err := s.SeedDepartments(); err != nil {
		return err
	}
	if err := s.SeedClasses(); err != nil {
		return err
	}
	// Subjects require teachers, but teachers require departments/subjects?
	// Circular dependency?
	// Teachers table has `subjects_taught` (ID array) or `department_id`.
	// Subjects table has `teacher_id`.
	// Let's create Teachers first, then Subjects?
	// But Teachers need Departments.
	// Order: Institutions -> Departments -> Classes -> Sections -> Teachers -> Subjects -> Students.

	// Re-ordered in logic below if necessary.
	// For now, let's keep basic structure and fix dependencies.

	// 3. User Roles (Admins, Accountants, Teachers, Students, Parents)
	if err := s.SeedUsers(); err != nil {
		return err
	}

	// 4. Subjects (after Teachers)
	if err := s.SeedSubjects(); err != nil {
		return err
	}

	logger.Info("Database seeding completed successfully")
	return nil
}

// SeedInstitutions creates default institutions
func (s *Seeder) SeedInstitutions() error {
	institutions := []models.Institution{
		{
			BaseModel:       models.BaseModel{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
			Name:            "Dhaka City College",
			Code:            "DCC",
			Address:         "Dhanmondi, Dhaka",
			Phone:           "+88029674115",
			Email:           "info@dhakacitycollege.edu.bd",
			PrincipalName:   "Prof. Bedar Uddin Ahmed",
			EstablishedYear: 1957,
			IsActive:        true,
		},
		{
			BaseModel:       models.BaseModel{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222")},
			Name:            "Test High School",
			Code:            "THS",
			Address:         "Test Street, Test City",
			Phone:           "+1000000000",
			Email:           "info@testschool.com",
			PrincipalName:   "Test Principal",
			EstablishedYear: 2020,
			IsActive:        true,
		},
	}

	for _, inst := range institutions {
		var count int64
		s.db.Model(&models.Institution{}).Where("code = ?", inst.Code).Count(&count)
		if count == 0 {
			// Ensure ID is set (using MustParse above) - Create with specific ID to help other seeders reference it
			if err := s.db.Create(&inst).Error; err != nil {
				return err
			}
			logger.Info("Institution seeded", zap.String("name", inst.Name))
		}
	}
	return nil
}

// SeedSuperAdmin creates a default super admin user if not exists
func (s *Seeder) SeedSuperAdmin() error {
	// Keep existing implementation but maybe update ID to be consistent?
	// User already exists check handles it.
	var count int64
	s.db.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&count)

	if count > 0 {
		return nil
	}

	// ... previous code for Super Admin ...
	// NOTE: This function's full body was not replaced, assuming I keep lines 36-86 as is?
	// The Prompt says "ReplacementContent" replaces from StartLine.
	// I should INCLUDE SeedSuperAdmin in ReplacementContent if I am replacing lines covering it.
	// Lines 24-123 covers `SeedAll`, `SeedSuperAdmin`, `SeedTestInstitution`.
	// I will include `SeedSuperAdmin` fully.

	hashedPassword, err := utils.HashPassword("Admin@123")
	if err != nil {
		return err
	}

	superAdmin := &models.User{
		BaseModel:    models.BaseModel{ID: uuid.New()},
		Email:        "superadmin@campus.local",
		PasswordHash: hashedPassword,
		Role:         models.RoleSuperAdmin,
		IsActive:     true,
	}

	if err := s.db.Create(superAdmin).Error; err != nil {
		return err
	}

	profile := &models.UserProfile{
		BaseModel: models.BaseModel{ID: uuid.New()},
		UserID:    superAdmin.ID,
		FirstName: "Super",
		LastName:  "Admin",
	}

	if err := s.db.Create(profile).Error; err != nil {
		return err
	}

	logger.Info("Super admin seeded")
	return nil
}
