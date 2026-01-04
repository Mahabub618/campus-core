package models

import (
	"time"

	"github.com/google/uuid"
)

// Student represents a student in the system
type Student struct {
	TenantBaseModel
	UserID        uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	ClassID       *uuid.UUID `gorm:"type:uuid" json:"class_id,omitempty"`
	SectionID     *uuid.UUID `gorm:"type:uuid" json:"section_id,omitempty"`
	RollNumber    int        `json:"roll_number,omitempty"`
	AdmissionDate *time.Time `json:"admission_date,omitempty"`
	BloodGroup    string     `gorm:"size:5" json:"blood_group,omitempty"`
	MedicalInfo   string     `gorm:"type:text" json:"medical_info,omitempty"`

	// Relations
	User    *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parents []*Parent `gorm:"many2many:parent_student_relations;" json:"parents,omitempty"`
}

// TableName specifies the table name for Student
func (Student) TableName() string {
	return "students"
}
