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
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// TeacherService handles teacher management logic
type TeacherService struct {
	repo       *repository.TeacherRepository
	userRepo   *repository.UserRepository
	db         *gorm.DB
	jwtManager *utils.JWTManager
}

func NewTeacherService(repo *repository.TeacherRepository, userRepo *repository.UserRepository, db *gorm.DB, jwtManager *utils.JWTManager) *TeacherService {
	return &TeacherService{
		repo:       repo,
		userRepo:   userRepo,
		db:         db,
		jwtManager: jwtManager,
	}
}

// CreateTeacher creates a new teacher
func (s *TeacherService) CreateTeacher(req *request.CreateTeacherRequest, creatorInstitutionID string) (*response.UserResponse, error) {
	// Validation
	if req.InstitutionID == "" {
		req.InstitutionID = creatorInstitutionID
	}
	if req.InstitutionID == "" {
		return nil, errors.New("institution_id is required")
	}

	// Password hashing
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	institutionID, _ := uuid.Parse(req.InstitutionID)

	// Create User & Teacher in transaction
	var teacherUser *models.User
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create User
		user := &models.User{
			BaseModel:    models.BaseModel{ID: uuid.New()},
			Email:        req.Email,
			Phone:        req.Phone,
			PasswordHash: hashedPassword,
			Role:         models.RoleTeacher,
			IsActive:     true,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 2. Create Profile
		profile := &models.UserProfile{
			BaseModel:     models.BaseModel{ID: uuid.New()},
			UserID:        user.ID,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			InstitutionID: &institutionID,
		}
		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		user.Profile = profile // For response
		teacherUser = user

		// 3. Create Teacher
		joiningDate, _ := time.Parse("2006-01-02", req.JoiningDate)
		var deptID *uuid.UUID
		if req.DepartmentID != "" {
			id, _ := uuid.Parse(req.DepartmentID)
			deptID = &id
		}

		teacher := &models.Teacher{
			TenantBaseModel: models.TenantBaseModel{
				BaseModel:     models.BaseModel{ID: uuid.New()},
				InstitutionID: institutionID,
			},
			UserID:         user.ID,
			JoiningDate:    &joiningDate,
			Qualifications: pq.StringArray(req.Qualifications),
			DepartmentID:   deptID,
		}
		if err := tx.Create(teacher).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	// Helper to convert to response (duplicate logic for now to avoid circular deps or complexity)
	resp := response.UserResponse{
		ID:       teacherUser.ID,
		Email:    teacherUser.Email,
		Phone:    teacherUser.Phone,
		Role:     teacherUser.Role,
		IsActive: teacherUser.IsActive,
		Profile: &response.ProfileResponse{
			ID:            teacherUser.Profile.ID,
			FirstName:     teacherUser.Profile.FirstName,
			LastName:      teacherUser.Profile.LastName,
			InstitutionID: teacherUser.Profile.InstitutionID,
		},
	}

	return &resp, nil
}

// GetAllTeachers returns all teachers for an institution
func (s *TeacherService) GetAllTeachers(institutionID string, params utils.PaginationParams) ([]response.UserResponse, utils.Pagination, error) {
	teachers, total, err := s.repo.FindAll(institutionID, params)
	if err != nil {
		return nil, utils.Pagination{}, utils.ErrInternalServer.Wrap(err)
	}

	var responses []response.UserResponse
	for _, t := range teachers {
		if t.User.ID != uuid.Nil {
			responses = append(responses, response.UserResponse{
				ID:       t.User.ID,
				Email:    t.User.Email,
				Phone:    t.User.Phone,
				Role:     t.User.Role,
				IsActive: t.User.IsActive,
				Profile: &response.ProfileResponse{
					ID:            t.User.Profile.ID,
					FirstName:     t.User.Profile.FirstName,
					LastName:      t.User.Profile.LastName,
					InstitutionID: t.User.Profile.InstitutionID,
				},
			})
		}
	}

	pagination := utils.NewPagination(params.Page, params.PerPage, total)
	return responses, pagination, nil
}

// GetTeacher gets a teacher by ID
func (s *TeacherService) GetTeacher(id uuid.UUID) (*response.UserResponse, error) {
	teacher, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := response.UserResponse{
		ID:       teacher.User.ID,
		Email:    teacher.User.Email,
		Phone:    teacher.User.Phone,
		Role:     teacher.User.Role,
		IsActive: teacher.User.IsActive,
		Profile: &response.ProfileResponse{
			ID:            teacher.User.Profile.ID,
			FirstName:     teacher.User.Profile.FirstName,
			LastName:      teacher.User.Profile.LastName,
			InstitutionID: teacher.User.Profile.InstitutionID,
		},
	}
	return &resp, nil
}

// UpdateTeacher updates a teacher
func (s *TeacherService) UpdateTeacher(id uuid.UUID, req *request.UpdateTeacherRequest, institutionID string) (*response.UserResponse, error) {
	teacher, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Verify tenant access
	if institutionID != "" && teacher.InstitutionID.String() != institutionID {
		return nil, utils.ErrCrossTenantAccess
	}

	// Update user fields
	if req.Email != "" && req.Email != teacher.User.Email {
		// Check email uniqueness
		var count int64
		if err := s.db.Model(&models.User{}).Where("email = ? AND id != ?", req.Email, teacher.User.ID).Count(&count).Error; err != nil {
			return nil, utils.ErrInternalServer.Wrap(err)
		}
		if count > 0 {
			return nil, utils.ErrEmailAlreadyExists
		}
		teacher.User.Email = req.Email
	}

	if req.Phone != "" {
		teacher.User.Phone = req.Phone
	}

	if req.IsActive != nil {
		teacher.User.IsActive = *req.IsActive
	}

	// Update profile fields
	if teacher.User.Profile != nil {
		if req.FirstName != "" {
			teacher.User.Profile.FirstName = req.FirstName
		}
		if req.LastName != "" {
			teacher.User.Profile.LastName = req.LastName
		}
	}

	// Update teacher-specific fields
	if req.Qualifications != nil {
		teacher.Qualifications = pq.StringArray(req.Qualifications)
	}

	if req.DepartmentID != "" {
		deptID, _ := uuid.Parse(req.DepartmentID)
		teacher.DepartmentID = &deptID
	}

	// Save changes in transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(teacher.User).Error; err != nil {
			return err
		}
		if teacher.User.Profile != nil {
			if err := tx.Save(teacher.User.Profile).Error; err != nil {
				return err
			}
		}
		if err := tx.Save(teacher).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, utils.ErrInternalServer.Wrap(err)
	}

	resp := response.UserResponse{
		ID:       teacher.User.ID,
		Email:    teacher.User.Email,
		Phone:    teacher.User.Phone,
		Role:     teacher.User.Role,
		IsActive: teacher.User.IsActive,
		Profile: &response.ProfileResponse{
			ID:            teacher.User.Profile.ID,
			FirstName:     teacher.User.Profile.FirstName,
			LastName:      teacher.User.Profile.LastName,
			InstitutionID: teacher.User.Profile.InstitutionID,
		},
	}
	return &resp, nil
}

// GetTeacherClasses gets a teacher's assigned classes
func (s *TeacherService) GetTeacherClasses(id uuid.UUID) ([]interface{}, error) {
	// Verify teacher exists
	if _, err := s.repo.FindByID(id); err != nil {
		return nil, err
	}

	// TODO: Implement when class_teacher_assignments table is available in Phase 3
	// For now, return empty array
	return []interface{}{}, nil
}

// GetTeacherSubjects gets a teacher's assigned subjects
func (s *TeacherService) GetTeacherSubjects(id uuid.UUID) ([]interface{}, error) {
	// Verify teacher exists
	if _, err := s.repo.FindByID(id); err != nil {
		return nil, err
	}

	// TODO: Implement when teacher_subject_assignments table is available in Phase 3
	// For now, return empty array
	return []interface{}{}, nil
}
