package models

import (
	"github.com/google/uuid"
)

// Department represents an academic department (e.g., Science, Arts)
type Department struct {
	BaseModel
	InstitutionID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"institution_id"`
	Name               string     `gorm:"size:100;not null" json:"name"`
	HeadOfDepartmentID *uuid.UUID `gorm:"type:uuid" json:"head_of_department_id,omitempty"`
	Description        string     `gorm:"type:text" json:"description,omitempty"`

	// Relations
	HeadOfDepartment *Teacher `gorm:"foreignKey:HeadOfDepartmentID" json:"head_of_department,omitempty"`
}

// TableName specifies the table name for Department
func (Department) TableName() string {
	return "departments"
}
