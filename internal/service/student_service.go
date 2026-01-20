package service

import (
	"errors"
	"time"

	"campus-core/internal/dto/request"
	"campus-core/internal/dto/response"
	"campus-core/internal/models"
	"campus-core/internal/repository"
	"campus-core/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudentService handles student management logic
type StudentService struct {
	repo       *repository.StudentRepository
	userRepo   *repository.UserRepository
	db         *gorm.DB
	jwtManager *utils.JWTManager
}

func NewStudentService(repo *repository.StudentRepository, userRepo *repository.UserRepository, db *gorm.DB, jwtManager *utils.JWTManager) *StudentService {
	return &StudentService{
		repo:       repo,
		userRepo:   userRepo,
		db:         db,
		jwtManager: jwtManager,
	}
}

// CreateStudent creates a new student
func (s *StudentService) CreateStudent(req *request.CreateStudentRequest, creatorInstitutionID string) (*response.UserResponse, error) {
	if req.InstitutionID == "" {
		req.InstitutionID = creatorInstitutionID
	}
	if req.InstitutionID == "" {
		return nil, errors.New("institution_id is required")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	institutionID, _ := uuid.Parse(req.InstitutionID)

	var studentUser *models.User
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create User
		user := &models.User{
			BaseModel:    models.BaseModel{ID: uuid.New()},
			Email:        req.Email,
			Phone:        req.Phone,
			PasswordHash: hashedPassword,
			Role:         models.RoleStudent,
			IsActive:     true,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 2. Create Profile
		profile := &models.UserProfile{
			BaseModel:       models.BaseModel{ID: uuid.New()},
			UserID:          user.ID,
			FirstName:       req.FirstName,
			LastName:        req.LastName,
			InstitutionID:   &institutionID,
			AdmissionNumber: req.AdmissionNumber,
			// DateOfBirth setting if available in request? Assuming DTO might need update or handled separately
		}
		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		user.Profile = profile
		studentUser = user

		// 3. Create Student
		admissionDate, _ := time.Parse("2006-01-02", req.AdmissionDate)
		var classID, sectionID *uuid.UUID
		if req.ClassID != "" {
			id, _ := uuid.Parse(req.ClassID)
			classID = &id
		}
		if req.SectionID != "" {
			id, _ := uuid.Parse(req.SectionID)
			sectionID = &id
		}

		student := &models.Student{
			TenantBaseModel: models.TenantBaseModel{
				BaseModel:     models.BaseModel{ID: uuid.New()},
				InstitutionID: institutionID,
			},
			UserID:        user.ID,
			AdmissionDate: &admissionDate,
			RollNumber:    req.RollNumber,
			ClassID:       classID,
			SectionID:     sectionID,
			BloodGroup:    req.BloodGroup,
			MedicalInfo:   req.MedicalInfo,
		}
		if err := tx.Create(student).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	resp := response.UserResponse{
		ID:       studentUser.ID,
		Email:    studentUser.Email,
		Phone:    studentUser.Phone,
		Role:     studentUser.Role,
		IsActive: studentUser.IsActive,
		Profile: &response.ProfileResponse{
			ID:            studentUser.Profile.ID,
			FirstName:     studentUser.Profile.FirstName,
			LastName:      studentUser.Profile.LastName,
			InstitutionID: studentUser.Profile.InstitutionID,
		},
	}

	return &resp, nil
}

// GetAllStudents returns all students
func (s *StudentService) GetAllStudents(institutionID string, params utils.PaginationParams) ([]response.UserResponse, utils.Pagination, error) {
	students, total, err := s.repo.FindAll(institutionID, "", "", params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.UserResponse
	for _, st := range students {
		if st.User.ID != uuid.Nil {
			responses = append(responses, response.UserResponse{
				ID:       st.User.ID,
				Email:    st.User.Email,
				Phone:    st.User.Phone,
				Role:     st.User.Role,
				IsActive: st.User.IsActive,
				Profile: &response.ProfileResponse{
					ID:            st.User.Profile.ID,
					FirstName:     st.User.Profile.FirstName,
					LastName:      st.User.Profile.LastName,
					InstitutionID: st.User.Profile.InstitutionID,
				},
			})
		}
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return responses, pagination, nil
}

// GetStudent gets a student by ID
func (s *StudentService) GetStudent(id uuid.UUID) (*response.UserResponse, error) {
	student, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := response.UserResponse{
		ID:       student.User.ID,
		Email:    student.User.Email,
		Phone:    student.User.Phone,
		Role:     student.User.Role,
		IsActive: student.User.IsActive,
		Profile: &response.ProfileResponse{
			ID:            student.User.Profile.ID,
			FirstName:     student.User.Profile.FirstName,
			LastName:      student.User.Profile.LastName,
			InstitutionID: student.User.Profile.InstitutionID,
		},
	}
	return &resp, nil
}

// UpdateStudent updates a student
func (s *StudentService) UpdateStudent(id uuid.UUID, req *request.UpdateStudentRequest, institutionID string) (*response.UserResponse, error) {
	student, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Verify tenant access
	if institutionID != "" && student.InstitutionID.String() != institutionID {
		return nil, utils.ErrCrossTenantAccess
	}

	// Update user fields
	if req.Email != "" && req.Email != student.User.Email {
		var count int64
		if err := s.db.Model(&models.User{}).Where("email = ? AND id != ?", req.Email, student.User.ID).Count(&count).Error; err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if count > 0 {
			return nil, utils.ErrEmailAlreadyExists
		}
		student.User.Email = req.Email
	}

	if req.Phone != "" {
		student.User.Phone = req.Phone
	}

	if req.IsActive != nil {
		student.User.IsActive = *req.IsActive
	}

	// Update profile fields
	if student.User.Profile != nil {
		if req.FirstName != "" {
			student.User.Profile.FirstName = req.FirstName
		}
		if req.LastName != "" {
			student.User.Profile.LastName = req.LastName
		}
	}

	// Update student-specific fields
	if req.ClassID != "" {
		classID, _ := uuid.Parse(req.ClassID)
		student.ClassID = &classID
	}

	if req.SectionID != "" {
		sectionID, _ := uuid.Parse(req.SectionID)
		student.SectionID = &sectionID
	}

	if req.RollNumber != nil {
		student.RollNumber = *req.RollNumber
	}

	if req.BloodGroup != "" {
		student.BloodGroup = req.BloodGroup
	}

	if req.MedicalInfo != "" {
		student.MedicalInfo = req.MedicalInfo
	}

	// Save changes in transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(student.User).Error; err != nil {
			return err
		}
		if student.User.Profile != nil {
			if err := tx.Save(student.User.Profile).Error; err != nil {
				return err
			}
		}
		if err := tx.Save(student).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	resp := response.UserResponse{
		ID:       student.User.ID,
		Email:    student.User.Email,
		Phone:    student.User.Phone,
		Role:     student.User.Role,
		IsActive: student.User.IsActive,
		Profile: &response.ProfileResponse{
			ID:            student.User.Profile.ID,
			FirstName:     student.User.Profile.FirstName,
			LastName:      student.User.Profile.LastName,
			InstitutionID: student.User.Profile.InstitutionID,
		},
	}
	return &resp, nil
}

// GetStudentParents gets a student's linked parents
func (s *StudentService) GetStudentParents(id uuid.UUID) ([]response.ParentRelationResponse, error) {
	student, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Load relations
	var relations []models.ParentStudentRelation
	if err := s.db.Preload("Parent.User.Profile").Where("student_id = ?", student.ID).Find(&relations).Error; err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.ParentRelationResponse
	for _, rel := range relations {
		if rel.Parent != nil && rel.Parent.User != nil {
			responses = append(responses, response.ParentRelationResponse{
				ParentID:     rel.ParentID,
				Relationship: rel.Relationship,
				IsPrimary:    rel.IsPrimary,
				Parent: response.UserResponse{
					ID:       rel.Parent.User.ID,
					Email:    rel.Parent.User.Email,
					Phone:    rel.Parent.User.Phone,
					Role:     rel.Parent.User.Role,
					IsActive: rel.Parent.User.IsActive,
					Profile: &response.ProfileResponse{
						ID:        rel.Parent.User.Profile.ID,
						FirstName: rel.Parent.User.Profile.FirstName,
						LastName:  rel.Parent.User.Profile.LastName,
					},
				},
			})
		}
	}

	return responses, nil
}

// LinkParent links a parent to a student
func (s *StudentService) LinkParent(studentID uuid.UUID, req *request.LinkParentRequest) error {
	// Verify student exists
	student, err := s.repo.FindByID(studentID)
	if err != nil {
		return err
	}

	parentID, err := uuid.Parse(req.ParentID)
	if err != nil {
		return utils.ErrInvalidUUID
	}

	// Verify parent exists and belongs to same institution
	var parent models.Parent
	if err := s.db.Where("id = ? AND institution_id = ?", parentID, student.InstitutionID).First(&parent).Error; err != nil {
		return utils.ErrInvalidParentStudentLink
	}

	// Check if relation already exists
	var count int64
	if err := s.db.Model(&models.ParentStudentRelation{}).
		Where("parent_id = ? AND student_id = ?", parentID, studentID).
		Count(&count).Error; err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}
	if count > 0 {
		return utils.ErrResourceExists
	}

	// Create relation
	relation := &models.ParentStudentRelation{
		BaseModel:    models.BaseModel{ID: uuid.New()},
		ParentID:     parentID,
		StudentID:    studentID,
		Relationship: req.Relationship,
		IsPrimary:    req.IsPrimary,
	}

	if err := s.db.Create(relation).Error; err != nil {
		return utils.ErrInternalServer.Wrap(err)
	}

	return nil
}

// UnlinkParent removes a parent-student relationship
func (s *StudentService) UnlinkParent(studentID, parentID uuid.UUID) error {
	// Verify student exists
	if _, err := s.repo.FindByID(studentID); err != nil {
		return err
	}

	// Delete the relation
	result := s.db.Where("parent_id = ? AND student_id = ?", parentID, studentID).
		Delete(&models.ParentStudentRelation{})

	if result.Error != nil {
		return utils.ErrInternalServer.Wrap(result.Error)
	}

	if result.RowsAffected == 0 {
		return utils.ErrResourceNotFound
	}

	return nil
}
