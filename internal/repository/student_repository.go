package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudentRepository handles student data
type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *StudentRepository) FindByID(id uuid.UUID) (*models.Student, error) {
	var student models.Student
	if err := r.db.Preload("User.Profile").First(&student, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) FindByUserID(userID uuid.UUID) (*models.Student, error) {
	var student models.Student
	if err := r.db.Preload("User.Profile").First(&student, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) Update(student *models.Student) error {
	return r.db.Save(student).Error
}

func (r *StudentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Student{}, "id = ?", id).Error
}

// FindAll returns filtered students (class, section filters can be added)
func (r *StudentRepository) FindAll(institutionID string, classID, sectionID string, params utils.PaginationParams) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	db := r.db.Model(&models.Student{}).Preload("User.Profile")

	if institutionID != "" {
		db = db.Where("institution_id = ?", institutionID)
	}
	if classID != "" {
		db = db.Where("class_id = ?", classID)
	}
	if sectionID != "" {
		db = db.Where("section_id = ?", sectionID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(utils.Paginate(params)).Find(&students).Error; err != nil {
		return nil, 0, err
	}

	return students, total, nil
}
