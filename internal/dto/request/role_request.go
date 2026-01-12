package request

// CreateTeacherRequest represents a request to create a teacher
type CreateTeacherRequest struct {
	RegisterRequest
	Qualifications []string `json:"qualifications"`
	JoiningDate    string   `json:"joining_date" binding:"required,datetime=2006-01-02"`
	DepartmentID   string   `json:"department_id" binding:"omitempty,uuid"`
}

// CreateStudentRequest represents a request to create a student
type CreateStudentRequest struct {
	RegisterRequest
	AdmissionNumber string `json:"admission_number" binding:"required"`
	AdmissionDate   string `json:"admission_date" binding:"required,datetime=2006-01-02"`
	RollNumber      int    `json:"roll_number"`
	ClassID         string `json:"class_id" binding:"omitempty,uuid"`
	SectionID       string `json:"section_id" binding:"omitempty,uuid"`
	BloodGroup      string `json:"blood_group"`
	MedicalInfo     string `json:"medical_info"`
}

// CreateParentRequest represents a request to create a parent
type CreateParentRequest struct {
	RegisterRequest
	Occupation       string `json:"occupation"`
	OfficeAddress    string `json:"office_address"`
	EmergencyContact string `json:"emergency_contact"`
}

// CreateAccountantRequest represents a request to create an accountant
type CreateAccountantRequest struct {
	RegisterRequest
	Qualification string `json:"qualification"`
	JoiningDate   string `json:"joining_date" binding:"required,datetime=2006-01-02"`
}
