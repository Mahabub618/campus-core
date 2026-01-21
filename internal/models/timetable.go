package models

import (
	"github.com/google/uuid"
)

// DayOfWeek represents days of the week
type DayOfWeek string

const (
	Sunday    DayOfWeek = "SUNDAY"
	Monday    DayOfWeek = "MONDAY"
	Tuesday   DayOfWeek = "TUESDAY"
	Wednesday DayOfWeek = "WEDNESDAY"
	Thursday  DayOfWeek = "THURSDAY"
	Friday    DayOfWeek = "FRIDAY"
	Saturday  DayOfWeek = "SATURDAY"
)

// Timetable represents a scheduled class period
type Timetable struct {
	BaseModel
	InstitutionID  uuid.UUID `gorm:"type:uuid;not null;index" json:"institution_id"`
	AcademicYearID uuid.UUID `gorm:"type:uuid;not null;index" json:"academic_year_id"`
	ClassID        uuid.UUID `gorm:"type:uuid;not null;index" json:"class_id"`
	SectionID      uuid.UUID `gorm:"type:uuid;not null;index" json:"section_id"`
	SubjectID      uuid.UUID `gorm:"type:uuid;not null;index" json:"subject_id"`
	TeacherID      uuid.UUID `gorm:"type:uuid;not null;index" json:"teacher_id"`
	DayOfWeek      DayOfWeek `gorm:"size:20;not null" json:"day_of_week"`
	StartTime      string    `gorm:"size:10;not null" json:"start_time"` // Format: "09:00"
	EndTime        string    `gorm:"size:10;not null" json:"end_time"`   // Format: "09:45"
	RoomNumber     string    `gorm:"size:50" json:"room_number,omitempty"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`

	// Relations
	Institution  *Institution  `gorm:"foreignKey:InstitutionID" json:"institution,omitempty"`
	AcademicYear *AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academic_year,omitempty"`
	Class        *Class        `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Section      *Section      `gorm:"foreignKey:SectionID" json:"section,omitempty"`
	Subject      *Subject      `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Teacher      *Teacher      `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

// TableName specifies the table name for Timetable
func (Timetable) TableName() string {
	return "timetables"
}

// Period represents a time slot in the school day
type Period struct {
	BaseModel
	InstitutionID uuid.UUID `gorm:"type:uuid;not null;index" json:"institution_id"`
	Name          string    `gorm:"size:50;not null" json:"name"` // e.g., "Period 1", "Lunch Break"
	StartTime     string    `gorm:"size:10;not null" json:"start_time"`
	EndTime       string    `gorm:"size:10;not null" json:"end_time"`
	Order         int       `gorm:"not null" json:"order"` // Order in the day
	IsBreak       bool      `gorm:"default:false" json:"is_break"`

	// Relations
	Institution *Institution `gorm:"foreignKey:InstitutionID" json:"institution,omitempty"`
}

// TableName specifies the table name for Period
func (Period) TableName() string {
	return "periods"
}
