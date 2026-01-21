package models

import (
	"time"

	"github.com/google/uuid"
)

// AcademicYear represents an academic year in the system
type AcademicYear struct {
	BaseModel
	InstitutionID uuid.UUID `gorm:"type:uuid;not null;index" json:"institution_id"`
	Name          string    `gorm:"size:50;not null" json:"name"` // e.g., "2025-2026"
	StartDate     time.Time `gorm:"not null" json:"start_date"`
	EndDate       time.Time `gorm:"not null" json:"end_date"`
	IsCurrent     bool      `gorm:"default:false" json:"is_current"`
	Description   string    `gorm:"type:text" json:"description,omitempty"`

	// Relations
	Institution *Institution `gorm:"foreignKey:InstitutionID" json:"institution,omitempty"`
}

// TableName specifies the table name for AcademicYear
func (AcademicYear) TableName() string {
	return "academic_years"
}

// Term represents a term/semester within an academic year
type Term struct {
	BaseModel
	AcademicYearID uuid.UUID `gorm:"type:uuid;not null;index" json:"academic_year_id"`
	Name           string    `gorm:"size:50;not null" json:"name"` // e.g., "First Term", "Second Semester"
	StartDate      time.Time `gorm:"not null" json:"start_date"`
	EndDate        time.Time `gorm:"not null" json:"end_date"`
	IsCurrent      bool      `gorm:"default:false" json:"is_current"`

	// Relations
	AcademicYear *AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academic_year,omitempty"`
}

// TableName specifies the table name for Term
func (Term) TableName() string {
	return "terms"
}
