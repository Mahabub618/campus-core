package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SectionFilter holds filter criteria for sections
type SectionFilter struct {
	ClassID string
	Search  string
}

// SectionRepository handles database operations for sections
type SectionRepository struct {
	db *gorm.DB
}

// NewSectionRepository creates a new section repository
func NewSectionRepository(db *gorm.DB) *SectionRepository {
	return &SectionRepository{db: db}
}

// FindByID finds a section by ID
func (r *SectionRepository) FindByID(id uuid.UUID) (*models.Section, error) {
	var section models.Section
	err := r.db.Preload("Class").First(&section, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &section, nil
}

// FindByClassID finds all sections for a class
func (r *SectionRepository) FindByClassID(classID uuid.UUID) ([]models.Section, error) {
	var sections []models.Section
	err := r.db.Where("class_id = ?", classID).Order("name ASC").Find(&sections).Error
	return sections, err
}

// FindAll finds all sections with filters
func (r *SectionRepository) FindAll(filter SectionFilter, params utils.PaginationParams) ([]models.Section, int64, error) {
	var sections []models.Section
	var total int64

	query := r.db.Model(&models.Section{})

	// Apply filters
	if filter.ClassID != "" {
		query = query.Where("class_id = ?", filter.ClassID)
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
	err := query.Preload("Class").
		Order("name ASC").Offset(offset).Limit(params.PerPage).Find(&sections).Error
	if err != nil {
		return nil, 0, err
	}

	return sections, total, nil
}

// Create creates a new section
func (r *SectionRepository) Create(section *models.Section) error {
	return r.db.Create(section).Error
}

// Update updates a section
func (r *SectionRepository) Update(section *models.Section) error {
	return r.db.Save(section).Error
}

// Delete soft deletes a section
func (r *SectionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Section{}, "id = ?", id).Error
}

// NameExistsInClass checks if a section name exists for a class
func (r *SectionRepository) NameExistsInClass(name string, classID uuid.UUID, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.Section{}).
		Where("name = ? AND class_id = ?", name, classID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// GetSectionStudentCount gets the count of students in a section
func (r *SectionRepository) GetSectionStudentCount(sectionID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Student{}).Where("section_id = ?", sectionID).Count(&count).Error
	return count, err
}

// GetSectionStudents gets all students in a section
func (r *SectionRepository) GetSectionStudents(sectionID uuid.UUID) ([]models.Student, error) {
	var students []models.Student
	err := r.db.Where("section_id = ?", sectionID).
		Preload("User").Preload("User.Profile").
		Find(&students).Error
	return students, err
}
