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
