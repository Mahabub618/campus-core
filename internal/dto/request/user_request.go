package request

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"omitempty"`
	FirstName string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName  string `json:"last_name" binding:"omitempty,min=1,max=100"`
	IsActive  *bool  `json:"is_active" binding:"omitempty"`
}

// UpdateProfileRequest represents a request to update user's own profile
type UpdateProfileRequest struct {
	FirstName       string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName        string `json:"last_name" binding:"omitempty,min=1,max=100"`
	DateOfBirth     string `json:"date_of_birth" binding:"omitempty,datetime=2006-01-02"`
	Gender          string `json:"gender" binding:"omitempty,oneof=male female other"`
	Address         string `json:"address" binding:"omitempty"`
	ProfileImageURL string `json:"profile_image_url" binding:"omitempty,url"`
}

// UpdateTeacherRequest represents a request to update a teacher
type UpdateTeacherRequest struct {
	Email          string   `json:"email" binding:"omitempty,email"`
	Phone          string   `json:"phone" binding:"omitempty"`
	FirstName      string   `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName       string   `json:"last_name" binding:"omitempty,min=1,max=100"`
	Qualifications []string `json:"qualifications" binding:"omitempty"`
	DepartmentID   string   `json:"department_id" binding:"omitempty,uuid"`
	IsActive       *bool    `json:"is_active" binding:"omitempty"`
}

// UpdateStudentRequest represents a request to update a student
type UpdateStudentRequest struct {
	Email       string `json:"email" binding:"omitempty,email"`
	Phone       string `json:"phone" binding:"omitempty"`
	FirstName   string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName    string `json:"last_name" binding:"omitempty,min=1,max=100"`
	ClassID     string `json:"class_id" binding:"omitempty,uuid"`
	SectionID   string `json:"section_id" binding:"omitempty,uuid"`
	RollNumber  *int   `json:"roll_number" binding:"omitempty"`
	BloodGroup  string `json:"blood_group" binding:"omitempty"`
	MedicalInfo string `json:"medical_info" binding:"omitempty"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// UpdateParentRequest represents a request to update a parent
type UpdateParentRequest struct {
	Email            string `json:"email" binding:"omitempty,email"`
	Phone            string `json:"phone" binding:"omitempty"`
	FirstName        string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName         string `json:"last_name" binding:"omitempty,min=1,max=100"`
	Occupation       string `json:"occupation" binding:"omitempty"`
	OfficeAddress    string `json:"office_address" binding:"omitempty"`
	EmergencyContact string `json:"emergency_contact" binding:"omitempty"`
	IsActive         *bool  `json:"is_active" binding:"omitempty"`
}

// UpdateAccountantRequest represents a request to update an accountant
type UpdateAccountantRequest struct {
	Email         string `json:"email" binding:"omitempty,email"`
	Phone         string `json:"phone" binding:"omitempty"`
	FirstName     string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName      string `json:"last_name" binding:"omitempty,min=1,max=100"`
	Qualification string `json:"qualification" binding:"omitempty"`
	IsActive      *bool  `json:"is_active" binding:"omitempty"`
}

// LinkParentRequest represents a request to link a parent to a student
type LinkParentRequest struct {
	ParentID     string `json:"parent_id" binding:"required,uuid"`
	Relationship string `json:"relationship" binding:"required,oneof=father mother guardian"`
	IsPrimary    bool   `json:"is_primary"`
}
