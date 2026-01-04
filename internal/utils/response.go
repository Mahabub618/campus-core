package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// Success sends a success response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithData sends a success response with data only
func SuccessWithData(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 created response
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// OK sends a 200 OK response
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

// NoContent sends a 204 no content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, data interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, err error) {
	response := ErrorResponse{
		Success: false,
		Error:   err.Error(),
	}

	// Check if it's an AppError to get more details
	if appErr, ok := err.(*AppError); ok {
		response.Code = appErr.Code
		response.Details = appErr.Details
		statusCode = appErr.StatusCode
	}

	c.JSON(statusCode, response)
}

// ErrorWithCode sends an error response with a specific code
func ErrorWithCode(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, details map[string]string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error:   "Validation failed",
		Code:    "VAL_001",
		Details: details,
	})
}

// BadRequest sends a 400 bad request response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// Unauthorized sends a 401 unauthorized response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "AUTH_004",
	})
}

// Forbidden sends a 403 forbidden response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "AUTHZ_001",
	})
}

// NotFound sends a 404 not found response
func NotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Success: false,
		Error:   resource + " not found",
		Code:    "RES_001",
	})
}

// InternalServerError sends a 500 internal server error response
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "SYS_001",
	})
}
