package models

import (
	"time"

	"github.com/google/uuid"
)

// Accountant represents an accountant in the system
type Accountant struct {
	TenantBaseModel

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user,omitempty"`

	JoiningDate   *time.Time `json:"joining_date"`
	Qualification string     `gorm:"type:varchar(255)" json:"qualification"`
}

// TableName specifies the table name
func (Accountant) TableName() string {
	return "accountants"
}
