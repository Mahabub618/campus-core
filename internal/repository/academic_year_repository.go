package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AcademicYearFilter holds filter criteria for academic years
type AcademicYearFilter struct {
	InstitutionID string
	IsCurrent     *bool
	Search        string
}

// AcademicYearRepository handles database operations for academic years
type AcademicYearRepository struct {
	db *gorm.DB
}

// NewAcademicYearRepository creates a new academic year repository
func NewAcademicYearRepository(db *gorm.DB) *AcademicYearRepository {
	return &AcademicYearRepository{db: db}
}

// FindByID finds an academic year by ID
func (r *AcademicYearRepository) FindByID(id uuid.UUID) (*models.AcademicYear, error) {
	var ay models.AcademicYear
	err := r.db.First(&ay, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &ay, nil
}

// FindByIDWithInstitution finds an academic year by ID with institution filter
func (r *AcademicYearRepository) FindByIDWithInstitution(id, institutionID uuid.UUID) (*models.AcademicYear, error) {
	var ay models.AcademicYear
	err := r.db.First(&ay, "id = ? AND institution_id = ?", id, institutionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &ay, nil
}

// FindAll finds all academic years with filters
func (r *AcademicYearRepository) FindAll(filter AcademicYearFilter, params utils.PaginationParams) ([]models.AcademicYear, int64, error) {
	var academicYears []models.AcademicYear
	var total int64

	query := r.db.Model(&models.AcademicYear{})

	// Apply filters
	if filter.InstitutionID != "" {
		query = query.Where("institution_id = ?", filter.InstitutionID)
	}
	if filter.IsCurrent != nil {
		query = query.Where("is_current = ?", *filter.IsCurrent)
	}
	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (params.Page - 1) * params.PerPage
	err := query.Order("start_date DESC").Offset(offset).Limit(params.PerPage).Find(&academicYears).Error
	if err != nil {
		return nil, 0, err
	}

	return academicYears, total, nil
}

// FindCurrent finds the current academic year for an institution
func (r *AcademicYearRepository) FindCurrent(institutionID uuid.UUID) (*models.AcademicYear, error) {
	var ay models.AcademicYear
	err := r.db.First(&ay, "institution_id = ? AND is_current = ?", institutionID, true).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &ay, nil
}

// Create creates a new academic year
func (r *AcademicYearRepository) Create(ay *models.AcademicYear) error {
	return r.db.Create(ay).Error
}

// Update updates an academic year
func (r *AcademicYearRepository) Update(ay *models.AcademicYear) error {
	return r.db.Save(ay).Error
}

// Delete soft deletes an academic year
func (r *AcademicYearRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.AcademicYear{}, "id = ?", id).Error
}

// SetCurrent sets an academic year as current and unsets others
func (r *AcademicYearRepository) SetCurrent(id, institutionID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Unset current for all academic years in the institution
		if err := tx.Model(&models.AcademicYear{}).
			Where("institution_id = ? AND is_current = ?", institutionID, true).
			Update("is_current", false).Error; err != nil {
			return err
		}
		// Set the specified academic year as current
		return tx.Model(&models.AcademicYear{}).
			Where("id = ? AND institution_id = ?", id, institutionID).
			Update("is_current", true).Error
	})
}

// NameExists checks if an academic year name exists for an institution
func (r *AcademicYearRepository) NameExists(name string, institutionID uuid.UUID, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.AcademicYear{}).
		Where("name = ? AND institution_id = ?", name, institutionID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
