package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubjectFilter holds filter criteria for subjects
type SubjectFilter struct {
	InstitutionID string
	ClassID       string
	TeacherID     string
	IsElective    *bool
	Search        string
}

// SubjectRepository handles database operations for subjects
type SubjectRepository struct {
	db *gorm.DB
}

// NewSubjectRepository creates a new subject repository
func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

// FindByID finds a subject by ID
func (r *SubjectRepository) FindByID(id uuid.UUID) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.Preload("Class").Preload("Teacher").First(&subject, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &subject, nil
}

// FindByIDWithInstitution finds a subject by ID with institution filter
func (r *SubjectRepository) FindByIDWithInstitution(id, institutionID uuid.UUID) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.Preload("Class").Preload("Teacher").
		First(&subject, "id = ? AND institution_id = ?", id, institutionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &subject, nil
}

// FindAll finds all subjects with filters
func (r *SubjectRepository) FindAll(filter SubjectFilter, params utils.PaginationParams) ([]models.Subject, int64, error) {
	var subjects []models.Subject
	var total int64

	query := r.db.Model(&models.Subject{})

	// Apply filters
	if filter.InstitutionID != "" {
		query = query.Where("institution_id = ?", filter.InstitutionID)
	}
	if filter.ClassID != "" {
		query = query.Where("class_id = ?", filter.ClassID)
	}
	if filter.TeacherID != "" {
		query = query.Where("teacher_id = ?", filter.TeacherID)
	}
	if filter.IsElective != nil {
		query = query.Where("is_elective = ?", *filter.IsElective)
	}
	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR code ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (params.Page - 1) * params.PerPage
	err := query.Preload("Class").Preload("Teacher").
		Order("name ASC").Offset(offset).Limit(params.PerPage).Find(&subjects).Error
	if err != nil {
		return nil, 0, err
	}

	return subjects, total, nil
}

// FindByClassID finds all subjects for a class
func (r *SubjectRepository) FindByClassID(classID uuid.UUID) ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Where("class_id = ?", classID).
		Preload("Teacher").
		Order("name ASC").Find(&subjects).Error
	return subjects, err
}

// FindByTeacherID finds all subjects assigned to a teacher
func (r *SubjectRepository) FindByTeacherID(teacherID uuid.UUID) ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Where("teacher_id = ?", teacherID).
		Preload("Class").
		Order("name ASC").Find(&subjects).Error
	return subjects, err
}

// Create creates a new subject
func (r *SubjectRepository) Create(subject *models.Subject) error {
	return r.db.Create(subject).Error
}

// Update updates a subject
func (r *SubjectRepository) Update(subject *models.Subject) error {
	return r.db.Save(subject).Error
}

// Delete soft deletes a subject
func (r *SubjectRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Subject{}, "id = ?", id).Error
}

// NameExistsInClass checks if a subject name exists for a class
func (r *SubjectRepository) NameExistsInClass(name string, classID uuid.UUID, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.Subject{}).
		Where("name = ? AND class_id = ?", name, classID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// CodeExists checks if a subject code exists for an institution
func (r *SubjectRepository) CodeExists(code string, institutionID uuid.UUID, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.Subject{}).
		Where("code = ? AND institution_id = ?", code, institutionID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// AssignTeacher assigns a teacher to a subject
func (r *SubjectRepository) AssignTeacher(subjectID, teacherID uuid.UUID) error {
	return r.db.Model(&models.Subject{}).
		Where("id = ?", subjectID).
		Update("teacher_id", teacherID).Error
}

// UnassignTeacher removes teacher assignment from a subject
func (r *SubjectRepository) UnassignTeacher(subjectID uuid.UUID) error {
	return r.db.Model(&models.Subject{}).
		Where("id = ?", subjectID).
		Update("teacher_id", nil).Error
}
