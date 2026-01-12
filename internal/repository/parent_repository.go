package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ParentRepository handles parent data
type ParentRepository struct {
	db *gorm.DB
}

func NewParentRepository(db *gorm.DB) *ParentRepository {
	return &ParentRepository{db: db}
}

func (r *ParentRepository) Create(parent *models.Parent) error {
	return r.db.Create(parent).Error
}

func (r *ParentRepository) FindByID(id uuid.UUID) (*models.Parent, error) {
	var parent models.Parent
	if err := r.db.Preload("User.Profile").First(&parent, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}
	return &parent, nil
}

func (r *ParentRepository) FindAll(institutionID string, params utils.PaginationParams) ([]models.Parent, int64, error) {
	var parents []models.Parent
	var total int64

	db := r.db.Model(&models.Parent{}).Preload("User.Profile")

	if institutionID != "" {
		db = db.Where("institution_id = ?", institutionID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(utils.Paginate(params)).Find(&parents).Error; err != nil {
		return nil, 0, err
	}

	return parents, total, nil
}

func (r *ParentRepository) Update(parent *models.Parent) error {
	return r.db.Save(parent).Error
}

func (r *ParentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Parent{}, "id = ?", id).Error
}
