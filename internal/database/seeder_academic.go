package database

import (
	"campus-core/internal/models"

	"campus-core/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SeedDepartments creates default departments
func (s *Seeder) SeedDepartments() error {
	institutions := []string{"DCC", "THS"}

	deptNames := []string{"Science", "Commerce", "Humanities"}

	for _, code := range institutions {
		var inst models.Institution
		if err := s.db.Where("code = ?", code).First(&inst).Error; err != nil {
			logger.Error("Institution not found for seeding departments", zap.String("code", code))
			continue
		}

		for _, name := range deptNames {
			var count int64
			s.db.Model(&models.Department{}).Where("institution_id = ? AND name = ?", inst.ID, name).Count(&count)
			if count == 0 {
				dept := &models.Department{
					BaseModel: models.BaseModel{ID: uuid.New()},
					TenantBaseModel: models.TenantBaseModel{
						InstitutionID: inst.ID,
					},
					Name:        name,
					Description: name + " Department",
				}
				if err := s.db.Create(dept).Error; err != nil {
					return err
				}
				logger.Info("Department seeded", zap.String("name", name), zap.String("institution", code))
			}
		}
	}
	return nil
}

// SeedClasses creates default classes and sections
func (s *Seeder) SeedClasses() error {
	// Define classes per institution
	classMap := map[string][]string{
		"DCC": {"Class 11", "Class 12"},
		"THS": {"Class 6", "Class 7", "Class 8", "Class 9", "Class 10"},
	}

	for code, classes := range classMap {
		var inst models.Institution
		if err := s.db.Where("code = ?", code).First(&inst).Error; err != nil {
			continue
		}

		for _, className := range classes {
			var class models.Class
			err := s.db.Where("institution_id = ? AND name = ?", inst.ID, className).First(&class).Error

			if err != nil {
				// Create class
				class = models.Class{
					BaseModel: models.BaseModel{ID: uuid.New()},
					TenantBaseModel: models.TenantBaseModel{
						InstitutionID: inst.ID,
					},
					Name:         className,
					SectionCount: 2,
					Capacity:     50,
				}
				if err := s.db.Create(&class).Error; err != nil {
					return err
				}
				logger.Info("Class seeded", zap.String("name", className), zap.String("institution", code))
			}

			// Create Sections (A, B)
			sections := []string{"A", "B"}
			for _, secName := range sections {
				var count int64
				s.db.Model(&models.Section{}).Where("class_id = ? AND name = ?", class.ID, secName).Count(&count)
				if count == 0 {
					section := &models.Section{
						BaseModel:  models.BaseModel{ID: uuid.New()},
						ClassID:    class.ID,
						Name:       secName,
						RoomNumber: "101", // Dummy
						Capacity:   40,
					}
					if err := s.db.Create(section).Error; err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// SeedSubjects creates default subjects
func (s *Seeder) SeedSubjects() error {
	// Subjects per Institution Type or just generic
	// We'll assign some common subjects to all classes

	commonSubjects := []string{"Bangla", "English", "Mathematics"}
	scienceSubjects := []string{"Physics", "Chemistry", "Biology"}

	// Get all classes
	var classes []models.Class
	if err := s.db.Find(&classes).Error; err != nil {
		return err
	}

	for _, class := range classes {
		// Fetch institution to check type/code if needed, but for now just add common subjects
		// Add Common Subjects
		for _, subName := range commonSubjects {
			s.createSubjectIfNotExists(class.ID, class.InstitutionID, subName, false)
		}

		// Add Science subjects only for higher classes (just dummy logic: Class 9, 10, 11, 12)
		if class.Name == "Class 9" || class.Name == "Class 10" || class.Name == "Class 11" || class.Name == "Class 12" {
			for _, subName := range scienceSubjects {
				s.createSubjectIfNotExists(class.ID, class.InstitutionID, subName, true) // Elective? or not
			}
		}
	}
	return nil
}

func (s *Seeder) createSubjectIfNotExists(classID uuid.UUID, institutionID uuid.UUID, name string, isElective bool) {
	var count int64
	s.db.Model(&models.Subject{}).Where("class_id = ? AND name = ?", classID, name).Count(&count)
	if count == 0 {
		subject := &models.Subject{
			BaseModel: models.BaseModel{ID: uuid.New()},
			TenantBaseModel: models.TenantBaseModel{
				InstitutionID: institutionID,
			},
			ClassID:     &classID,
			Name:        name,
			Code:        name[0:3] + "-101", // Dummy code
			IsElective:  isElective,
			CreditHours: 3.0,
		}
		s.db.Create(subject)
		logger.Info("Subject seeded", zap.String("name", name), zap.String("class_id", classID.String()))
	}
}
