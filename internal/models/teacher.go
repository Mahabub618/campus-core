package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Teacher represents a teacher in the system
type Teacher struct {
	TenantBaseModel
	UserID         uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Qualifications pq.StringArray `gorm:"type:text[]" json:"qualifications,omitempty"`
	JoiningDate    *time.Time     `json:"joining_date,omitempty"`
	DepartmentID   *uuid.UUID     `gorm:"type:uuid" json:"department_id,omitempty"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Teacher
func (Teacher) TableName() string {
	return "teachers"
}
