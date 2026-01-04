package models

import (
	"github.com/google/uuid"
)

// Institution represents a school/institution in the system
type Institution struct {
	BaseModel
	Name            string `gorm:"size:255;not null" json:"name"`
	Code            string `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Address         string `gorm:"type:text" json:"address,omitempty"`
	Phone           string `gorm:"size:20" json:"phone,omitempty"`
	Email           string `gorm:"size:255" json:"email,omitempty"`
	PrincipalName   string `gorm:"size:255" json:"principal_name,omitempty"`
	EstablishedYear int    `json:"established_year,omitempty"`
	LogoURL         string `gorm:"size:500" json:"logo_url,omitempty"`
	AcademicYear    string `gorm:"size:20" json:"academic_year,omitempty"`
	IsActive        bool   `gorm:"default:true" json:"is_active"`
}

// TableName specifies the table name for Institution
func (Institution) TableName() string {
	return "institutions"
}

// InstitutionStats represents statistics for an institution
type InstitutionStats struct {
	TotalStudents int64     `json:"total_students"`
	TotalTeachers int64     `json:"total_teachers"`
	TotalParents  int64     `json:"total_parents"`
	TotalClasses  int64     `json:"total_classes"`
	ActiveUsers   int64     `json:"active_users"`
	InstitutionID uuid.UUID `json:"-"`
}
