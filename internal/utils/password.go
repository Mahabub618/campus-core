package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the default bcrypt cost factor
	DefaultCost = 12
)

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength checks if password meets minimum requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case 'a' <= c && c <= 'z':
			hasLower = true
		case '0' <= c && c <= '9':
			hasDigit = true
		case c == '@' || c == '#' || c == '$' || c == '%' || c == '!' || c == '&' || c == '*':
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return ErrPasswordRequirements
	}

	return nil
}
