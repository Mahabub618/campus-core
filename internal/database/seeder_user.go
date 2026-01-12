package database

import (
	"fmt"
	"time"

	"campus-core/internal/models"
	"campus-core/internal/utils"
	"campus-core/pkg/logger"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// SeedUsers creates default users for roles
func (s *Seeder) SeedUsers() error {
	institutions := []string{"DCC", "THS"}

	for _, code := range institutions {
		var inst models.Institution
		if err := s.db.Where("code = ?", code).First(&inst).Error; err != nil {
			continue
		}

		// 1. Admin
		if err := s.seedRoleUser(inst.ID, "Admin", "User", "admin@"+inst.Code+".edu.bd", models.RoleAdmin); err != nil {
			return err
		}

		// 2. Accountant
		if err := s.seedAccountant(inst.ID); err != nil {
			return err
		}

		// 3. Teachers
		if err := s.seedTeachers(inst.ID); err != nil {
			return err
		}

		// 4. Students & Parents
		if err := s.seedStudents(inst.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) seedRoleUser(institutionID uuid.UUID, firstName, lastName, email, role string) error {
	var count int64
	s.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, _ := utils.HashPassword("Pass@123")
	user := &models.User{
		BaseModel:    models.BaseModel{ID: uuid.New()},
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
		IsActive:     true,
	}
	if err := s.db.Create(user).Error; err != nil {
		return err
	}

	profile := &models.UserProfile{
		BaseModel:     models.BaseModel{ID: uuid.New()},
		UserID:        user.ID,
		InstitutionID: &institutionID,
		FirstName:     firstName,
		LastName:      lastName,
	}
	if err := s.db.Create(profile).Error; err != nil {
		return err
	}
	logger.Info("User seeded", zap.String("email", email), zap.String("role", role))
	return nil
}

func (s *Seeder) seedAccountant(institutionID uuid.UUID) error {
	email := "accountant@" + s.getInstCode(institutionID) + ".edu.bd"
	if err := s.seedRoleUser(institutionID, "Accountant", "Staff", email, models.RoleAccountant); err != nil {
		return err
	}
	// Create Accountant specific record
	var user models.User
	s.db.Where("email = ?", email).First(&user)

	var count int64
	s.db.Model(&models.Accountant{}).Where("user_id = ?", user.ID).Count(&count)
	if count == 0 {
		joinDate := time.Now()
		acc := &models.Accountant{
			TenantBaseModel: models.TenantBaseModel{
				BaseModel:     models.BaseModel{ID: uuid.New()},
				InstitutionID: institutionID,
			},
			UserID:        user.ID,
			JoiningDate:   &joinDate,
			Qualification: "BBA",
		}
		s.db.Create(acc)
	}
	return nil
}

func (s *Seeder) seedTeachers(institutionID uuid.UUID) error {
	// Create 5 teachers
	// First get a department to assign
	var dept models.Department
	s.db.Where("institution_id = ?", institutionID).First(&dept)

	for i := 1; i <= 5; i++ {
		email := fmt.Sprintf("teacher%d@%s.edu.bd", i, s.getInstCode(institutionID))
		if err := s.seedRoleUser(institutionID, "Teacher", fmt.Sprintf("%d", i), email, models.RoleTeacher); err != nil {
			return err
		}

		var user models.User
		s.db.Where("email = ?", email).First(&user)

		var count int64
		s.db.Model(&models.Teacher{}).Where("user_id = ?", user.ID).Count(&count)

		if count == 0 {
			joinDate := time.Now()
			teacher := &models.Teacher{
				TenantBaseModel: models.TenantBaseModel{
					BaseModel:     models.BaseModel{ID: uuid.New()},
					InstitutionID: institutionID,
				},
				UserID:         user.ID,
				JoiningDate:    &joinDate,
				Qualifications: pq.StringArray{"M.Sc", "B.Ed"},
				DepartmentID:   &dept.ID, // Assign to first department found
			}
			s.db.Create(teacher)
		}
	}
	return nil
}

func (s *Seeder) seedStudents(institutionID uuid.UUID) error {
	// Create 10 students
	// Assign to a class/section
	var section models.Section
	// Join with Class to ensure institution match
	// GORM join complicated, just get class first
	var class models.Class
	s.db.Where("institution_id = ?", institutionID).First(&class)
	if class.ID != uuid.Nil {
		s.db.Where("class_id = ?", class.ID).First(&section)
	}

	for i := 1; i <= 10; i++ {
		email := fmt.Sprintf("student%d@%s.edu.bd", i, s.getInstCode(institutionID))
		if err := s.seedRoleUser(institutionID, "Student", fmt.Sprintf("%d", i), email, models.RoleStudent); err != nil {
			return err
		}

		var user models.User
		s.db.Where("email = ?", email).First(&user)

		var count int64
		s.db.Model(&models.Student{}).Where("user_id = ?", user.ID).Count(&count)

		if count == 0 {
			admDate := time.Now()
			roll := i
			student := &models.Student{
				TenantBaseModel: models.TenantBaseModel{
					BaseModel:     models.BaseModel{ID: uuid.New()},
					InstitutionID: institutionID,
				},
				UserID:        user.ID,
				AdmissionDate: &admDate,
				RollNumber:    roll,
				ClassID:       &class.ID,
				SectionID:     &section.ID,
				BloodGroup:    "B+",
			}
			if err := s.db.Create(student).Error; err != nil {
				logger.Error("Failed to create student", zap.Error(err))
			}

			// Seed Parent for this student
			parentEmail := fmt.Sprintf("parent%d@%s.edu.bd", i, s.getInstCode(institutionID))
			s.seedRoleUser(institutionID, "Parent", fmt.Sprintf("%d", i), parentEmail, models.RoleParent)

			var parentUser models.User
			s.db.Where("email = ?", parentEmail).First(&parentUser)

			// Create Parent Record
			var pCount int64
			s.db.Model(&models.Parent{}).Where("user_id = ?", parentUser.ID).Count(&pCount)
			var parentID uuid.UUID
			if pCount == 0 {
				parent := &models.Parent{
					TenantBaseModel: models.TenantBaseModel{
						BaseModel:     models.BaseModel{ID: uuid.New()},
						InstitutionID: institutionID,
					},
					UserID:     parentUser.ID,
					Occupation: "Business",
				}
				s.db.Create(parent)
				parentID = parent.ID
			} else {
				var existingParent models.Parent
				s.db.Where("user_id = ?", parentUser.ID).First(&existingParent)
				parentID = existingParent.ID
			}

			// Link Parent to Student
			if student.ID != uuid.Nil && parentID != uuid.Nil {
				relation := &models.ParentStudentRelation{
					BaseModel:    models.BaseModel{ID: uuid.New()},
					ParentID:     parentID,
					StudentID:    student.ID,
					Relationship: "Father",
					IsPrimary:    true,
				}
				s.db.Create(relation)
			}
		}
	}
	return nil
}

func (s *Seeder) getInstCode(id uuid.UUID) string {
	var inst models.Institution
	s.db.First(&inst, id)
	return inst.Code
}
