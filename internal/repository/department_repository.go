package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DepartmentFilter holds filter criteria for departments
type DepartmentFilter struct {
	InstitutionID string
	Search        string
}

// DepartmentRepository handles database operations for departments
type DepartmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository creates a new department repository
func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

// FindByID finds a department by ID
func (r *DepartmentRepository) FindByID(id uuid.UUID) (*models.Department, error) {
	var dept models.Department
	err := r.db.Preload("HeadOfDepartment").First(&dept, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &dept, nil
}

// FindByIDWithInstitution finds a department by ID with institution filter
func (r *DepartmentRepository) FindByIDWithInstitution(id, institutionID uuid.UUID) (*models.Department, error) {
	var dept models.Department
	err := r.db.Preload("HeadOfDepartment").
		First(&dept, "id = ? AND institution_id = ?", id, institutionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &dept, nil
}

// FindAll finds all departments with filters
func (r *DepartmentRepository) FindAll(filter DepartmentFilter, params utils.PaginationParams) ([]models.Department, int64, error) {
	var departments []models.Department
	var total int64

	query := r.db.Model(&models.Department{})

	// Apply filters
	if filter.InstitutionID != "" {
		query = query.Where("institution_id = ?", filter.InstitutionID)
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
	err := query.Preload("HeadOfDepartment").
		Order("name ASC").Offset(offset).Limit(params.PerPage).Find(&departments).Error
	if err != nil {
		return nil, 0, err
	}

	return departments, total, nil
}

// Create creates a new department
func (r *DepartmentRepository) Create(dept *models.Department) error {
	return r.db.Create(dept).Error
}

// Update updates a department
func (r *DepartmentRepository) Update(dept *models.Department) error {
	return r.db.Save(dept).Error
}

// Delete soft deletes a department
func (r *DepartmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Department{}, "id = ?", id).Error
}

// NameExists checks if a department name exists for an institution
func (r *DepartmentRepository) NameExists(name string, institutionID uuid.UUID, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.Department{}).
		Where("name = ? AND institution_id = ?", name, institutionID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// GetDepartmentStaff gets all teachers in a department
func (r *DepartmentRepository) GetDepartmentStaff(departmentID uuid.UUID) ([]models.Teacher, error) {
	var teachers []models.Teacher
	err := r.db.Where("department_id = ?", departmentID).
		Preload("User").Preload("User.Profile").
		Find(&teachers).Error
	return teachers, err
}

// GetStaffCount gets the count of staff in a department
func (r *DepartmentRepository) GetStaffCount(departmentID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Teacher{}).Where("department_id = ?", departmentID).Count(&count).Error
	return count, err
}
