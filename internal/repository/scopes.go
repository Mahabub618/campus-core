package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TenantScope filters queries by institution_id
func TenantScope(institutionID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if institutionID == "" {
			return db
		}
		// Parse UUID to ensure safety, though it should be validated by middleware
		if _, err := uuid.Parse(institutionID); err != nil {
			return db
		}
		return db.Where("institution_id = ?", institutionID)
	}
}

// ActiveScope filters by is_active = true
func ActiveScope(db *gorm.DB) *gorm.DB {
	return db.Where("is_active = ?", true)
}
