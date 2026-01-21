package response

import (
	"time"

	"github.com/google/uuid"
)

// AcademicYearResponse represents the response for an academic year
type AcademicYearResponse struct {
	ID            uuid.UUID `json:"id"`
	InstitutionID uuid.UUID `json:"institution_id"`
	Name          string    `json:"name"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsCurrent     bool      `json:"is_current"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ClassResponse represents the response for a class
type ClassResponse struct {
	ID             uuid.UUID         `json:"id"`
	InstitutionID  uuid.UUID         `json:"institution_id"`
	Name           string            `json:"name"`
	SectionCount   int               `json:"section_count"`
	ClassTeacherID *uuid.UUID        `json:"class_teacher_id,omitempty"`
	ClassTeacher   *TeacherBrief     `json:"class_teacher,omitempty"`
	Capacity       int               `json:"capacity,omitempty"`
	Sections       []SectionResponse `json:"sections,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// ClassBrief represents a brief class response (for nested objects)
type ClassBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// SectionResponse represents the response for a section
type SectionResponse struct {
	ID         uuid.UUID   `json:"id"`
	ClassID    uuid.UUID   `json:"class_id"`
	Name       string      `json:"name"`
	RoomNumber string      `json:"room_number,omitempty"`
	Capacity   int         `json:"capacity,omitempty"`
	Class      *ClassBrief `json:"class,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// SectionBrief represents a brief section response (for nested objects)
type SectionBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// SubjectResponse represents the response for a subject
type SubjectResponse struct {
	ID            uuid.UUID     `json:"id"`
	InstitutionID uuid.UUID     `json:"institution_id"`
	ClassID       *uuid.UUID    `json:"class_id,omitempty"`
	TeacherID     *uuid.UUID    `json:"teacher_id,omitempty"`
	Name          string        `json:"name"`
	Code          string        `json:"code,omitempty"`
	IsElective    bool          `json:"is_elective"`
	CreditHours   float64       `json:"credit_hours,omitempty"`
	Class         *ClassBrief   `json:"class,omitempty"`
	Teacher       *TeacherBrief `json:"teacher,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// SubjectBrief represents a brief subject response (for nested objects)
type SubjectBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code,omitempty"`
}

// TeacherBrief represents a brief teacher response (for nested objects)
type TeacherBrief struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

// DepartmentResponse represents the response for a department
type DepartmentResponse struct {
	ID                 uuid.UUID     `json:"id"`
	InstitutionID      uuid.UUID     `json:"institution_id"`
	Name               string        `json:"name"`
	HeadOfDepartmentID *uuid.UUID    `json:"head_of_department_id,omitempty"`
	HeadOfDepartment   *TeacherBrief `json:"head_of_department,omitempty"`
	Description        string        `json:"description,omitempty"`
	StaffCount         int64         `json:"staff_count,omitempty"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}

// TimetableResponse represents the response for a timetable entry
type TimetableResponse struct {
	ID             uuid.UUID     `json:"id"`
	InstitutionID  uuid.UUID     `json:"institution_id"`
	AcademicYearID uuid.UUID     `json:"academic_year_id"`
	ClassID        uuid.UUID     `json:"class_id"`
	SectionID      uuid.UUID     `json:"section_id"`
	SubjectID      uuid.UUID     `json:"subject_id"`
	TeacherID      uuid.UUID     `json:"teacher_id"`
	DayOfWeek      string        `json:"day_of_week"`
	StartTime      string        `json:"start_time"`
	EndTime        string        `json:"end_time"`
	RoomNumber     string        `json:"room_number,omitempty"`
	IsActive       bool          `json:"is_active"`
	Class          *ClassBrief   `json:"class,omitempty"`
	Section        *SectionBrief `json:"section,omitempty"`
	Subject        *SubjectBrief `json:"subject,omitempty"`
	Teacher        *TeacherBrief `json:"teacher,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// DayTimetable represents timetable entries grouped by day
type DayTimetable struct {
	Day     string              `json:"day"`
	Entries []TimetableResponse `json:"entries"`
}

// WeekTimetableResponse represents a full week's timetable
type WeekTimetableResponse struct {
	Days []DayTimetable `json:"days"`
}
