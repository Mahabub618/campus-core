package request

import "time"

// CreateAcademicYearRequest represents the request to create an academic year
type CreateAcademicYearRequest struct {
	Name        string    `json:"name" binding:"required,min=3,max=50"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required,gtfield=StartDate"`
	IsCurrent   bool      `json:"is_current"`
	Description string    `json:"description" binding:"max=500"`
}

// UpdateAcademicYearRequest represents the request to update an academic year
type UpdateAcademicYearRequest struct {
	Name        string     `json:"name" binding:"omitempty,min=3,max=50"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	IsCurrent   *bool      `json:"is_current"`
	Description string     `json:"description" binding:"max=500"`
}

// CreateClassRequest represents the request to create a class
type CreateClassRequest struct {
	Name           string `json:"name" binding:"required,min=1,max=50"`
	ClassTeacherID string `json:"class_teacher_id" binding:"omitempty,uuid"`
	Capacity       int    `json:"capacity" binding:"omitempty,min=1,max=500"`
}

// UpdateClassRequest represents the request to update a class
type UpdateClassRequest struct {
	Name           string `json:"name" binding:"omitempty,min=1,max=50"`
	ClassTeacherID string `json:"class_teacher_id" binding:"omitempty,uuid"`
	Capacity       *int   `json:"capacity" binding:"omitempty,min=1,max=500"`
}

// CreateSectionRequest represents the request to create a section
type CreateSectionRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=50"`
	RoomNumber string `json:"room_number" binding:"max=20"`
	Capacity   int    `json:"capacity" binding:"omitempty,min=1,max=100"`
}

// UpdateSectionRequest represents the request to update a section
type UpdateSectionRequest struct {
	Name       string `json:"name" binding:"omitempty,min=1,max=50"`
	RoomNumber string `json:"room_number" binding:"max=20"`
	Capacity   *int   `json:"capacity" binding:"omitempty,min=1,max=100"`
}

// CreateSubjectRequest represents the request to create a subject
type CreateSubjectRequest struct {
	ClassID     string  `json:"class_id" binding:"omitempty,uuid"`
	TeacherID   string  `json:"teacher_id" binding:"omitempty,uuid"`
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Code        string  `json:"code" binding:"omitempty,max=20"`
	IsElective  bool    `json:"is_elective"`
	CreditHours float64 `json:"credit_hours" binding:"omitempty,min=0,max=10"`
}

// UpdateSubjectRequest represents the request to update a subject
type UpdateSubjectRequest struct {
	ClassID     string   `json:"class_id" binding:"omitempty,uuid"`
	TeacherID   string   `json:"teacher_id" binding:"omitempty,uuid"`
	Name        string   `json:"name" binding:"omitempty,min=1,max=100"`
	Code        string   `json:"code" binding:"omitempty,max=20"`
	IsElective  *bool    `json:"is_elective"`
	CreditHours *float64 `json:"credit_hours" binding:"omitempty,min=0,max=10"`
}

// AssignTeacherRequest represents the request to assign a teacher to a subject
type AssignTeacherRequest struct {
	TeacherID string `json:"teacher_id" binding:"required,uuid"`
}

// CreateDepartmentRequest represents the request to create a department
type CreateDepartmentRequest struct {
	Name               string `json:"name" binding:"required,min=1,max=100"`
	HeadOfDepartmentID string `json:"head_of_department_id" binding:"omitempty,uuid"`
	Description        string `json:"description" binding:"max=500"`
}

// UpdateDepartmentRequest represents the request to update a department
type UpdateDepartmentRequest struct {
	Name               string `json:"name" binding:"omitempty,min=1,max=100"`
	HeadOfDepartmentID string `json:"head_of_department_id" binding:"omitempty,uuid"`
	Description        string `json:"description" binding:"max=500"`
}

// CreateTimetableRequest represents the request to create a timetable entry
type CreateTimetableRequest struct {
	AcademicYearID string `json:"academic_year_id" binding:"required,uuid"`
	ClassID        string `json:"class_id" binding:"required,uuid"`
	SectionID      string `json:"section_id" binding:"required,uuid"`
	SubjectID      string `json:"subject_id" binding:"required,uuid"`
	TeacherID      string `json:"teacher_id" binding:"required,uuid"`
	DayOfWeek      string `json:"day_of_week" binding:"required,oneof=SUNDAY MONDAY TUESDAY WEDNESDAY THURSDAY FRIDAY SATURDAY"`
	StartTime      string `json:"start_time" binding:"required"` // Format: "09:00"
	EndTime        string `json:"end_time" binding:"required"`   // Format: "09:45"
	RoomNumber     string `json:"room_number" binding:"max=50"`
}

// UpdateTimetableRequest represents the request to update a timetable entry
type UpdateTimetableRequest struct {
	AcademicYearID string `json:"academic_year_id" binding:"omitempty,uuid"`
	ClassID        string `json:"class_id" binding:"omitempty,uuid"`
	SectionID      string `json:"section_id" binding:"omitempty,uuid"`
	SubjectID      string `json:"subject_id" binding:"omitempty,uuid"`
	TeacherID      string `json:"teacher_id" binding:"omitempty,uuid"`
	DayOfWeek      string `json:"day_of_week" binding:"omitempty,oneof=SUNDAY MONDAY TUESDAY WEDNESDAY THURSDAY FRIDAY SATURDAY"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	RoomNumber     string `json:"room_number" binding:"max=50"`
	IsActive       *bool  `json:"is_active"`
}

// BulkTimetableRequest represents the request to create multiple timetable entries
type BulkTimetableRequest struct {
	Entries []CreateTimetableRequest `json:"entries" binding:"required,min=1,dive"`
}
