package utils

import (
	"math"

	"gorm.io/gorm"
)

// Pagination holds pagination information
type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// PaginationParams holds pagination request parameters
type PaginationParams struct {
	Page    int
	PerPage int
}

// DefaultPagination returns default pagination parameters
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:    1,
		PerPage: 20,
	}
}

// NewPaginationParams creates pagination params with validation
func NewPaginationParams(page, perPage int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100 // Max limit
	}
	return PaginationParams{
		Page:    page,
		PerPage: perPage,
	}
}

// GetOffset returns the offset for database queries
func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

// GetLimit returns the limit for database queries
func (p PaginationParams) GetLimit() int {
	return p.PerPage
}

// NewPagination creates a new Pagination object
func NewPagination(page, perPage int, totalItems int64) Pagination {
	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	if totalPages == 0 && totalItems > 0 {
		totalPages = 1
	}

	return Pagination{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}
}

// Paginate is a GORM scope function for pagination
func Paginate(params PaginationParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(params.GetOffset()).Limit(params.GetLimit())
	}
}

// CountAndPaginate counts total records and applies pagination
func CountAndPaginate(db *gorm.DB, model interface{}, params PaginationParams) (*gorm.DB, int64, error) {
	var totalItems int64

	// Count total items
	if err := db.Model(model).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	paginatedDB := db.Scopes(Paginate(params))

	return paginatedDB, totalItems, nil
}

// HasNextPage checks if there's a next page
func (p Pagination) HasNextPage() bool {
	return p.CurrentPage < p.TotalPages
}

// HasPreviousPage checks if there's a previous page
func (p Pagination) HasPreviousPage() bool {
	return p.CurrentPage > 1
}
