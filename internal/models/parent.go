package models

import (
	"github.com/google/uuid"
)

// Parent represents a parent/guardian in the system
type Parent struct {
	TenantBaseModel
	UserID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Occupation       string    `gorm:"size:100" json:"occupation,omitempty"`
	OfficeAddress    string    `gorm:"type:text" json:"office_address,omitempty"`
	EmergencyContact string    `gorm:"size:20" json:"emergency_contact,omitempty"`

	// Relations
	User     *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Children []*Student `gorm:"many2many:parent_student_relations;" json:"children,omitempty"`
}

// TableName specifies the table name for Parent
func (Parent) TableName() string {
	return "parents"
}

// ParentStudentRelation represents the relationship between parents and students
type ParentStudentRelation struct {
	BaseModel
	ParentID     uuid.UUID `gorm:"type:uuid;not null;index" json:"parent_id"`
	StudentID    uuid.UUID `gorm:"type:uuid;not null;index" json:"student_id"`
	Relationship string    `gorm:"size:50" json:"relationship"` // Father, Mother, Guardian
	IsPrimary    bool      `gorm:"default:false" json:"is_primary"`

	// Relations
	Parent  *Parent  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Student *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

// TableName specifies the table name for ParentStudentRelation
func (ParentStudentRelation) TableName() string {
	return "parent_student_relations"
}
