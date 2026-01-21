package models

import (
	"github.com/google/uuid"
)

// Class represents a student class (e.g., Class 10)
type Class struct {
	BaseModel
	InstitutionID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"institution_id"`
	Name           string     `gorm:"size:50;not null" json:"name"`
	SectionCount   int        `gorm:"default:1" json:"section_count"`
	ClassTeacherID *uuid.UUID `gorm:"type:uuid" json:"class_teacher_id,omitempty"`
	Capacity       int        `json:"capacity,omitempty"`

	// Relations
	ClassTeacher *Teacher  `gorm:"foreignKey:ClassTeacherID" json:"class_teacher,omitempty"`
	Sections     []Section `gorm:"foreignKey:ClassID" json:"sections,omitempty"`
}

// TableName specifies the table name for Class
func (Class) TableName() string {
	return "classes"
}

// Section represents a section within a class (e.g., Class 10 - Section A)
type Section struct {
	BaseModel
	ClassID    uuid.UUID `gorm:"type:uuid;not null" json:"class_id"`
	Name       string    `gorm:"size:50;not null" json:"name"`
	RoomNumber string    `gorm:"size:20" json:"room_number,omitempty"`
	Capacity   int       `json:"capacity,omitempty"`

	// Relations
	Class *Class `gorm:"foreignKey:ClassID" json:"class,omitempty"`
}

// TableName specifies the table name for Section
func (Section) TableName() string {
	return "sections"
}

// Subject represents an academic subject
type Subject struct {
	BaseModel
	InstitutionID uuid.UUID  `gorm:"type:uuid;not null;index" json:"institution_id"`
	ClassID       *uuid.UUID `gorm:"type:uuid" json:"class_id,omitempty"`
	TeacherID     *uuid.UUID `gorm:"type:uuid" json:"teacher_id,omitempty"`
	Name          string     `gorm:"size:100;not null" json:"name"`
	Code          string     `gorm:"size:20" json:"code,omitempty"`
	IsElective    bool       `gorm:"default:false" json:"is_elective"`
	CreditHours   float64    `gorm:"type:decimal(4,2)" json:"credit_hours,omitempty"`

	// Relations
	Class   *Class   `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Teacher *Teacher `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

// TableName specifies the table name for Subject
func (Subject) TableName() string {
	return "subjects"
}
