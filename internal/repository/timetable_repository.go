package repository

import (
	"errors"

	"campus-core/internal/models"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TimetableFilter holds filter criteria for timetable entries
type TimetableFilter struct {
	InstitutionID  string
	AcademicYearID string
	ClassID        string
	SectionID      string
	SubjectID      string
	TeacherID      string
	DayOfWeek      string
	IsActive       *bool
}

// TimetableRepository handles database operations for timetable
type TimetableRepository struct {
	db *gorm.DB
}

// NewTimetableRepository creates a new timetable repository
func NewTimetableRepository(db *gorm.DB) *TimetableRepository {
	return &TimetableRepository{db: db}
}

// FindByID finds a timetable entry by ID
func (r *TimetableRepository) FindByID(id uuid.UUID) (*models.Timetable, error) {
	var tt models.Timetable
	err := r.db.Preload("Class").Preload("Section").Preload("Subject").Preload("Teacher").
		First(&tt, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &tt, nil
}

// FindByIDWithInstitution finds a timetable entry by ID with institution filter
func (r *TimetableRepository) FindByIDWithInstitution(id, institutionID uuid.UUID) (*models.Timetable, error) {
	var tt models.Timetable
	err := r.db.Preload("Class").Preload("Section").Preload("Subject").Preload("Teacher").
		First(&tt, "id = ? AND institution_id = ?", id, institutionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &tt, nil
}

// FindAll finds all timetable entries with filters
func (r *TimetableRepository) FindAll(filter TimetableFilter, params utils.PaginationParams) ([]models.Timetable, int64, error) {
	var timetables []models.Timetable
	var total int64

	query := r.db.Model(&models.Timetable{})

	// Apply filters
	if filter.InstitutionID != "" {
		query = query.Where("institution_id = ?", filter.InstitutionID)
	}
	if filter.AcademicYearID != "" {
		query = query.Where("academic_year_id = ?", filter.AcademicYearID)
	}
	if filter.ClassID != "" {
		query = query.Where("class_id = ?", filter.ClassID)
	}
	if filter.SectionID != "" {
		query = query.Where("section_id = ?", filter.SectionID)
	}
	if filter.SubjectID != "" {
		query = query.Where("subject_id = ?", filter.SubjectID)
	}
	if filter.TeacherID != "" {
		query = query.Where("teacher_id = ?", filter.TeacherID)
	}
	if filter.DayOfWeek != "" {
		query = query.Where("day_of_week = ?", filter.DayOfWeek)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (params.Page - 1) * params.PerPage
	err := query.Preload("Class").Preload("Section").Preload("Subject").Preload("Teacher").
		Order("day_of_week ASC, start_time ASC").Offset(offset).Limit(params.PerPage).Find(&timetables).Error
	if err != nil {
		return nil, 0, err
	}

	return timetables, total, nil
}

// FindByClassID finds all timetable entries for a class
func (r *TimetableRepository) FindByClassID(classID uuid.UUID, academicYearID *uuid.UUID) ([]models.Timetable, error) {
	var timetables []models.Timetable
	query := r.db.Where("class_id = ? AND is_active = ?", classID, true)
	if academicYearID != nil {
		query = query.Where("academic_year_id = ?", *academicYearID)
	}
	err := query.Preload("Section").Preload("Subject").Preload("Teacher").
		Order("day_of_week ASC, start_time ASC").Find(&timetables).Error
	return timetables, err
}

// FindBySectionID finds all timetable entries for a section
func (r *TimetableRepository) FindBySectionID(sectionID uuid.UUID, academicYearID *uuid.UUID) ([]models.Timetable, error) {
	var timetables []models.Timetable
	query := r.db.Where("section_id = ? AND is_active = ?", sectionID, true)
	if academicYearID != nil {
		query = query.Where("academic_year_id = ?", *academicYearID)
	}
	err := query.Preload("Class").Preload("Subject").Preload("Teacher").
		Order("day_of_week ASC, start_time ASC").Find(&timetables).Error
	return timetables, err
}

// FindByTeacherID finds all timetable entries for a teacher
func (r *TimetableRepository) FindByTeacherID(teacherID uuid.UUID, academicYearID *uuid.UUID) ([]models.Timetable, error) {
	var timetables []models.Timetable
	query := r.db.Where("teacher_id = ? AND is_active = ?", teacherID, true)
	if academicYearID != nil {
		query = query.Where("academic_year_id = ?", *academicYearID)
	}
	err := query.Preload("Class").Preload("Section").Preload("Subject").
		Order("day_of_week ASC, start_time ASC").Find(&timetables).Error
	return timetables, err
}

// Create creates a new timetable entry
func (r *TimetableRepository) Create(tt *models.Timetable) error {
	return r.db.Create(tt).Error
}

// Update updates a timetable entry
func (r *TimetableRepository) Update(tt *models.Timetable) error {
	return r.db.Save(tt).Error
}

// Delete soft deletes a timetable entry
func (r *TimetableRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Timetable{}, "id = ?", id).Error
}

// CheckConflict checks for scheduling conflicts
// Returns true if there's a conflict
func (r *TimetableRepository) CheckConflict(tt *models.Timetable, excludeID *uuid.UUID) (bool, error) {
	var count int64

	// Check teacher conflict: same teacher, same day, overlapping time
	teacherQuery := r.db.Model(&models.Timetable{}).
		Where("teacher_id = ? AND day_of_week = ? AND is_active = ?", tt.TeacherID, tt.DayOfWeek, true).
		Where("((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
			tt.StartTime, tt.StartTime, tt.EndTime, tt.EndTime, tt.StartTime, tt.EndTime)
	if excludeID != nil {
		teacherQuery = teacherQuery.Where("id != ?", *excludeID)
	}
	if err := teacherQuery.Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	// Check section conflict: same section, same day, overlapping time
	sectionQuery := r.db.Model(&models.Timetable{}).
		Where("section_id = ? AND day_of_week = ? AND is_active = ?", tt.SectionID, tt.DayOfWeek, true).
		Where("((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
			tt.StartTime, tt.StartTime, tt.EndTime, tt.EndTime, tt.StartTime, tt.EndTime)
	if excludeID != nil {
		sectionQuery = sectionQuery.Where("id != ?", *excludeID)
	}
	if err := sectionQuery.Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	// Check room conflict if room is specified
	if tt.RoomNumber != "" {
		roomQuery := r.db.Model(&models.Timetable{}).
			Where("room_number = ? AND day_of_week = ? AND is_active = ?", tt.RoomNumber, tt.DayOfWeek, true).
			Where("((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
				tt.StartTime, tt.StartTime, tt.EndTime, tt.EndTime, tt.StartTime, tt.EndTime)
		if excludeID != nil {
			roomQuery = roomQuery.Where("id != ?", *excludeID)
		}
		if err := roomQuery.Count(&count).Error; err != nil {
			return false, err
		}
		if count > 0 {
			return true, nil
		}
	}

	return false, nil
}

// BulkCreate creates multiple timetable entries
func (r *TimetableRepository) BulkCreate(timetables []models.Timetable) error {
	return r.db.CreateInBatches(timetables, 100).Error
}

// DeleteByAcademicYear deletes all timetable entries for an academic year
func (r *TimetableRepository) DeleteByAcademicYear(academicYearID uuid.UUID) error {
	return r.db.Where("academic_year_id = ?", academicYearID).Delete(&models.Timetable{}).Error
}
