package models

import (
	"time"

	"github.com/google/uuid"
)

// Role constants
const (
	RoleSuperAdmin = "SUPER_ADMIN"
	RoleAdmin      = "ADMIN"
	RoleTeacher    = "TEACHER"
	RoleStudent    = "STUDENT"
	RoleParent     = "PARENT"
	RoleAccountant = "ACCOUNTANT"
)

// ValidRoles contains all valid user roles
var ValidRoles = []string{
	RoleSuperAdmin,
	RoleAdmin,
	RoleTeacher,
	RoleStudent,
	RoleParent,
	RoleAccountant,
}

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	for _, r := range ValidRoles {
		if r == role {
			return true
		}
	}
	return false
}

// User represents a user in the system
type User struct {
	BaseModel
	Email            string       `gorm:"size:255;uniqueIndex" json:"email,omitempty"`
	Phone            string       `gorm:"size:20" json:"phone,omitempty"`
	PasswordHash     string       `gorm:"size:255" json:"-"`
	Role             string       `gorm:"size:50;not null" json:"role"`
	IsActive         bool         `gorm:"default:true" json:"is_active"`
	LastLoginAt      *time.Time   `json:"last_login_at,omitempty"`
	RefreshToken     string       `gorm:"size:500" json:"-"`
	ResetToken       string       `gorm:"size:255" json:"-"`
	ResetTokenExpiry *time.Time   `json:"-"`
	Profile          *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// UserProfile represents profile information for a user
type UserProfile struct {
	BaseModel
	UserID          uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	InstitutionID   *uuid.UUID `gorm:"type:uuid;index" json:"institution_id,omitempty"`
	FirstName       string     `gorm:"size:100" json:"first_name"`
	LastName        string     `gorm:"size:100" json:"last_name"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"`
	Gender          string     `gorm:"size:10" json:"gender,omitempty"`
	Address         string     `gorm:"type:text" json:"address,omitempty"`
	ProfileImageURL string     `gorm:"size:500" json:"profile_image_url,omitempty"`
	EmployeeID      string     `gorm:"size:50" json:"employee_id,omitempty"`
	AdmissionNumber string     `gorm:"size:50" json:"admission_number,omitempty"`
	Occupation      string     `gorm:"size:100" json:"occupation,omitempty"`

	// Relations
	User        *User        `gorm:"foreignKey:UserID" json:"-"`
	Institution *Institution `gorm:"foreignKey:InstitutionID" json:"institution,omitempty"`
}

// TableName specifies the table name for UserProfile
func (UserProfile) TableName() string {
	return "user_profiles"
}

// FullName returns the user's full name
func (p *UserProfile) FullName() string {
	if p.FirstName == "" && p.LastName == "" {
		return ""
	}
	if p.LastName == "" {
		return p.FirstName
	}
	if p.FirstName == "" {
		return p.LastName
	}
	return p.FirstName + " " + p.LastName
}
