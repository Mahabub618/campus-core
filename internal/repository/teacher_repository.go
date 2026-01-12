package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TeacherRepository handles teacher data
type TeacherRepository struct {
	db *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) Create(teacher *models.Teacher) error {
	return r.db.Create(teacher).Error
}

func (r *TeacherRepository) FindByID(id uuid.UUID) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := r.db.Preload("User.Profile").First(&teacher, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}
	return &teacher, nil
}

func (r *TeacherRepository) FindByUserID(userID uuid.UUID) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := r.db.Preload("User.Profile").First(&teacher, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}
	return &teacher, nil
}

func (r *TeacherRepository) Update(teacher *models.Teacher) error {
	return r.db.Save(teacher).Error
}

func (r *TeacherRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Teacher{}, "id = ?", id).Error
}

func (r *TeacherRepository) FindAll(institutionID string, params utils.PaginationParams) ([]models.Teacher, int64, error) {
	var teachers []models.Teacher
	var total int64

	db := r.db.Model(&models.Teacher{}).Preload("User.Profile")

	if institutionID != "" {
		db = db.Where("institution_id = ?", institutionID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(utils.Paginate(params)).Find(&teachers).Error; err != nil {
		return nil, 0, err
	}

	return teachers, total, nil
}
