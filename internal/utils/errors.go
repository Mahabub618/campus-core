package utils

import (
	"fmt"
	"net/http"
)

// AppError represents a structured application error
type AppError struct {
	Code       string            `json:"code"`
	Message    string            `json:"error"`
	StatusCode int               `json:"-"`
	Details    map[string]string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewAppErrorWithDetails creates a new application error with details
func NewAppErrorWithDetails(code, message string, statusCode int, details map[string]string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    details,
	}
}

// Wrap wraps an error with additional context
func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    fmt.Sprintf("%s: %v", e.Message, err),
		StatusCode: e.StatusCode,
		Details:    e.Details,
	}
}

// Authentication Errors (AUTH_xxx)
var (
	ErrInvalidCredentials   = NewAppError("AUTH_001", "Invalid credentials", http.StatusUnauthorized)
	ErrTokenExpired         = NewAppError("AUTH_002", "Token has expired", http.StatusUnauthorized)
	ErrTokenInvalid         = NewAppError("AUTH_003", "Token is invalid", http.StatusUnauthorized)
	ErrTokenMissing         = NewAppError("AUTH_004", "Authorization token is required", http.StatusUnauthorized)
	ErrRefreshTokenExpired  = NewAppError("AUTH_005", "Refresh token has expired", http.StatusUnauthorized)
	ErrRefreshTokenInvalid  = NewAppError("AUTH_006", "Refresh token is invalid", http.StatusUnauthorized)
	ErrAccountDisabled      = NewAppError("AUTH_007", "Account is disabled", http.StatusForbidden)
	ErrAccountLocked        = NewAppError("AUTH_008", "Account is locked", http.StatusForbidden)
	ErrPasswordRequirements = NewAppError("AUTH_009", "Password does not meet requirements", http.StatusBadRequest)
	ErrResetTokenInvalid    = NewAppError("AUTH_010", "Password reset token is invalid", http.StatusBadRequest)
	ErrResetTokenExpired    = NewAppError("AUTH_011", "Password reset token has expired", http.StatusBadRequest)
	ErrTooManyLoginAttempts = NewAppError("AUTH_012", "Too many login attempts, please try again later", http.StatusTooManyRequests)
	ErrPasswordTooShort     = NewAppError("AUTH_009", "Password must be at least 8 characters", http.StatusBadRequest)
)

// Authorization Errors (AUTHZ_xxx)
var (
	ErrInsufficientPermissions = NewAppError("AUTHZ_001", "Insufficient permissions", http.StatusForbidden)
	ErrRoleNotAllowed          = NewAppError("AUTHZ_002", "Role not allowed for this action", http.StatusForbidden)
	ErrResourceAccessDenied    = NewAppError("AUTHZ_003", "Access to resource denied", http.StatusForbidden)
	ErrActionNotPermitted      = NewAppError("AUTHZ_004", "Action not permitted for your role", http.StatusForbidden)
	ErrCrossTenantAccess       = NewAppError("AUTHZ_005", "Cross-tenant access denied", http.StatusForbidden)
)

// Validation Errors (VAL_xxx)
var (
	ErrRequiredFieldMissing = NewAppError("VAL_001", "Required field missing", http.StatusBadRequest)
	ErrInvalidFieldFormat   = NewAppError("VAL_002", "Invalid field format", http.StatusBadRequest)
	ErrFieldOutOfRange      = NewAppError("VAL_003", "Field value out of range", http.StatusBadRequest)
	ErrInvalidDateFormat    = NewAppError("VAL_004", "Invalid date format", http.StatusBadRequest)
	ErrInvalidEmailFormat   = NewAppError("VAL_005", "Invalid email format", http.StatusBadRequest)
	ErrInvalidPhoneFormat   = NewAppError("VAL_006", "Invalid phone format", http.StatusBadRequest)
	ErrFieldTooLong         = NewAppError("VAL_007", "Field value too long", http.StatusBadRequest)
	ErrFieldTooShort        = NewAppError("VAL_008", "Field value too short", http.StatusBadRequest)
	ErrInvalidUUID          = NewAppError("VAL_009", "Invalid UUID format", http.StatusBadRequest)
	ErrInvalidEnumValue     = NewAppError("VAL_010", "Invalid enum value", http.StatusBadRequest)
	ErrUnprocessableEntity  = NewAppError("VAL_011", "Unprocessable entity", http.StatusUnprocessableEntity)
)

// Resource Errors (RES_xxx)
var (
	ErrResourceNotFound      = NewAppError("RES_001", "Resource not found", http.StatusNotFound)
	ErrNotFound              = ErrResourceNotFound // Alias for convenience
	ErrResourceExists        = NewAppError("RES_002", "Resource already exists", http.StatusConflict)
	ErrDuplicateEntry        = NewAppError("RES_003", "Duplicate entry", http.StatusConflict)
	ErrResourceInUse         = NewAppError("RES_004", "Resource is in use and cannot be deleted", http.StatusBadRequest)
	ErrResourceLimitExceeded = NewAppError("RES_005", "Resource limit exceeded", http.StatusBadRequest)
	ErrInvalidResourceState  = NewAppError("RES_006", "Invalid resource state", http.StatusBadRequest)
)

// User Management Errors (USER_xxx)
var (
	ErrUserNotFound              = NewAppError("USER_001", "User not found", http.StatusNotFound)
	ErrEmailAlreadyExists        = NewAppError("USER_002", "Email already registered", http.StatusConflict)
	ErrPhoneAlreadyExists        = NewAppError("USER_003", "Phone already registered", http.StatusConflict)
	ErrInvalidRoleAssignment     = NewAppError("USER_004", "Invalid role assignment", http.StatusBadRequest)
	ErrCannotDeleteSelf          = NewAppError("USER_005", "Cannot delete your own account", http.StatusBadRequest)
	ErrCannotDeactivateLastAdmin = NewAppError("USER_006", "Cannot deactivate the last admin", http.StatusBadRequest)
	ErrInvalidParentStudentLink  = NewAppError("USER_007", "Invalid parent-student link", http.StatusBadRequest)
)

// Institution Errors (INST_xxx)
var (
	ErrInstitutionNotFound   = NewAppError("INST_001", "Institution not found", http.StatusNotFound)
	ErrInstitutionCodeExists = NewAppError("INST_002", "Institution code already exists", http.StatusConflict)
	ErrInstitutionDisabled   = NewAppError("INST_003", "Institution is disabled", http.StatusBadRequest)
	ErrInstitutionIDRequired = NewAppError("INST_004", "X-Institution-ID header is required", http.StatusBadRequest)
	ErrUserNotInInstitution  = NewAppError("INST_005", "User does not belong to this institution", http.StatusForbidden)
)

// System Errors (SYS_xxx)
var (
	ErrInternalServer     = NewAppError("SYS_001", "Internal server error", http.StatusInternalServerError)
	ErrServiceUnavailable = NewAppError("SYS_002", "Service temporarily unavailable", http.StatusServiceUnavailable)
	ErrDatabaseError      = NewAppError("SYS_003", "Database error", http.StatusInternalServerError)
	ErrCacheError         = NewAppError("SYS_004", "Cache error", http.StatusInternalServerError)
	ErrRateLimitExceeded  = NewAppError("SYS_005", "Rate limit exceeded", http.StatusTooManyRequests)
	ErrWebSocketError     = NewAppError("SYS_006", "WebSocket connection error", http.StatusInternalServerError)
)
