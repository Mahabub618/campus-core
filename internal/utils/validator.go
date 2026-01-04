package utils

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// CustomValidator holds the custom validator instance
var CustomValidator *validator.Validate

// InitValidator initializes the custom validator with custom rules
func InitValidator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		CustomValidator = v

		// Register custom tag name function to use JSON tag names in error messages
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Register custom validators
		if err := v.RegisterValidation("role", validateRole); err != nil {
			return err
		}

		if err := v.RegisterValidation("phone", validatePhone); err != nil {
			return err
		}

		if err := v.RegisterValidation("password", validatePassword); err != nil {
			return err
		}
	}

	return nil
}

// validateRole validates that the role is one of the allowed values
func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	validRoles := []string{"SUPER_ADMIN", "ADMIN", "TEACHER", "STUDENT", "PARENT", "ACCOUNTANT"}

	for _, r := range validRoles {
		if r == role {
			return true
		}
	}
	return false
}

// validatePhone validates phone number format
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // Optional field
	}

	// Basic phone validation: allow +, digits, spaces, dashes
	phoneRegex := regexp.MustCompile(`^\+?[\d\s\-()]{10,20}$`)
	return phoneRegex.MatchString(phone)
}

// validatePassword validates password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, c := range password {
		switch {
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case 'a' <= c && c <= 'z':
			hasLower = true
		case '0' <= c && c <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

// FormatValidationErrors formats validation errors into a map
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "email":
				errors[field] = field + " must be a valid email address"
			case "min":
				errors[field] = field + " must be at least " + e.Param() + " characters"
			case "max":
				errors[field] = field + " must be at most " + e.Param() + " characters"
			case "role":
				errors[field] = field + " must be a valid role (SUPER_ADMIN, ADMIN, TEACHER, STUDENT, PARENT, ACCOUNTANT)"
			case "phone":
				errors[field] = field + " must be a valid phone number"
			case "password":
				errors[field] = field + " must be at least 8 characters with uppercase, lowercase, and digits"
			case "uuid":
				errors[field] = field + " must be a valid UUID"
			default:
				errors[field] = field + " is invalid"
			}
		}
	}

	return errors
}
