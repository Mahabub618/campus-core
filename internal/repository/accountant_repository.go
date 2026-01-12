package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountantRepository handles accountant data
type AccountantRepository struct {
	db *gorm.DB
}

func NewAccountantRepository(db *gorm.DB) *AccountantRepository {
	return &AccountantRepository{db: db}
}

func (r *AccountantRepository) Create(accountant *models.Accountant) error {
	return r.db.Create(accountant).Error
}

func (r *AccountantRepository) FindByID(id uuid.UUID) (*models.Accountant, error) {
	var accountant models.Accountant
	if err := r.db.Preload("User.Profile").First(&accountant, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, err
	}
	return &accountant, nil
}

func (r *AccountantRepository) FindAll(institutionID string, params utils.PaginationParams) ([]models.Accountant, int64, error) {
	var accountants []models.Accountant
	var total int64

	db := r.db.Model(&models.Accountant{}).Preload("User.Profile")

	if institutionID != "" {
		db = db.Where("institution_id = ?", institutionID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Scopes(utils.Paginate(params)).Find(&accountants).Error; err != nil {
		return nil, 0, err
	}

	return accountants, total, nil
}

func (r *AccountantRepository) Update(accountant *models.Accountant) error {
	return r.db.Save(accountant).Error
}

func (r *AccountantRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Accountant{}, "id = ?", id).Error
}
