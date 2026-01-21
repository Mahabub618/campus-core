package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ClassFilter holds filter criteria for classes
type ClassFilter struct {
	InstitutionID string
	Search        string
}

// ClassRepository handles database operations for classes
type ClassRepository struct {
	db *gorm.DB
}

// NewClassRepository creates a new class repository
func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

// FindByID finds a class by ID
func (r *ClassRepository) FindByID(id uuid.UUID) (*models.Class, error) {
	var class models.Class
	err := r.db.Preload("Sections").Preload("ClassTeacher").First(&class, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &class, nil
}

// FindByIDWithInstitution finds a class by ID with institution filter
func (r *ClassRepository) FindByIDWithInstitution(id, institutionID uuid.UUID) (*models.Class, error) {
	var class models.Class
	err := r.db.Preload("Sections").Preload("ClassTeacher").
		First(&class, "id = ? AND institution_id = ?", id, institutionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &class, nil
}

// FindAll finds all classes with filters
func (r *ClassRepository) FindAll(filter ClassFilter, params utils.PaginationParams) ([]models.Class, int64, error) {
	var classes []models.Class
	var total int64

	query := r.db.Model(&models.Class{})

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
	err := query.Preload("Sections").Preload("ClassTeacher").
		Order("name ASC").Offset(offset).Limit(params.PerPage).Find(&classes).Error
	if err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

// FindAllWithoutPagination finds all classes without pagination (for dropdowns)
func (r *ClassRepository) FindAllWithoutPagination(institutionID uuid.UUID) ([]models.Class, error) {
	var classes []models.Class
	err := r.db.Where("institution_id = ?", institutionID).Order("name ASC").Find(&classes).Error
	return classes, err
}

// Create creates a new class
func (r *ClassRepository) Create(class *models.Class) error {
	return r.db.Create(class).Error
}

// Update updates a class
func (r *ClassRepository) Update(class *models.Class) error {
	return r.db.Save(class).Error
}

// Delete soft deletes a class
func (r *ClassRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Class{}, "id = ?", id).Error
}

// NameExists checks if a class name exists for an institution
func (r *ClassRepository) NameExists(name string, institutionID uuid.UUID, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.Class{}).
		Where("name = ? AND institution_id = ?", name, institutionID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// GetClassStudentCount gets the count of students in a class
func (r *ClassRepository) GetClassStudentCount(classID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Student{}).Where("class_id = ?", classID).Count(&count).Error
	return count, err
}

// GetClassTeachers gets all teachers assigned to a class (via subjects or class teacher)
func (r *ClassRepository) GetClassTeachers(classID uuid.UUID) ([]models.Teacher, error) {
	var teachers []models.Teacher

	// Get class teacher and subject teachers
	err := r.db.Distinct().
		Joins("JOIN subjects ON subjects.teacher_id = teachers.id").
		Where("subjects.class_id = ?", classID).
		Or("teachers.id IN (SELECT class_teacher_id FROM classes WHERE id = ?)", classID).
		Preload("User").Preload("User.Profile").
		Find(&teachers).Error

	return teachers, err
}
