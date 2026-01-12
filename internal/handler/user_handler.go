package handler

import (
	"net/http"

	"campus-core/internal/dto/request"
	"campus-core/internal/middleware"
	"campus-core/internal/repository"
	"campus-core/internal/service"
	"campus-core/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user API requests
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser handles user creation (Admin)
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	creatorRole := middleware.GetUserRole(c)
	creatorInstID := middleware.GetInstitutionID(c)

	resp, err := h.service.CreateUser(&req, creatorRole, creatorInstID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.Created(c, "User created successfully", resp)
}

// GetAllUsers lists users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = utils.DefaultPagination()
	} else {
		params = utils.NewPaginationParams(params.Page, params.PerPage)
	}

	// Filters
	filter := repository.UserFilter{
		Role:          c.Query("role"),
		Search:        c.Query("search"),
		InstitutionID: middleware.GetInstitutionID(c), // Enforce tenant
	}

	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		filter.IsActive = &active
	}

	data, pagination, err := h.service.GetAllUsers(filter, params)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.Paginated(c, data, pagination)
}

// GetUser gets a single user
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	// Security: Check tenant access?
	// If FindByID doesn't filter by tenant, user from other tenant might be returned.
	// User Service/Repo should handle this via Tenant Scope if injected.
	// But `FindByID` in repo currently handles it? `FindByID` in `user_repository` doesn't use `institution_id` filter explicitly yet.
	// We need to ensure data privacy.
	// Ideally, `FindByID` should also check institution_id if the user is not Super Admin.
	// Or we rely on `TenantMiddleware` setting the scope, and repo using it.
	// My `FindByID` implementation in `user_repository` doesn't check scope. I should fix that.

	currentInstID := middleware.GetInstitutionID(c)
	if currentInstID != "" && user.Profile != nil && user.Profile.InstitutionID != nil && user.Profile.InstitutionID.String() != currentInstID {
		utils.Error(c, http.StatusNotFound, utils.ErrUserNotFound) // Pretend not found
		return
	}

	utils.OK(c, "", user)
}

// ToggleStatus updates user status
func (h *UserHandler) ToggleStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, utils.ErrInvalidUUID)
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid body")
		return
	}

	if err := h.service.ToggleStatus(id, req.IsActive); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.OK(c, "User status updated", nil)
}

// GetProfile gets current user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ErrTokenMissing)
		return
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.OK(c, "", user)
}

// UpdateProfile updates current user's profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ErrTokenMissing)
		return
	}

	var req struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	user, err := h.service.UpdateProfile(userID, req.FirstName, req.LastName)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.OK(c, "Profile updated successfully", user)
}
