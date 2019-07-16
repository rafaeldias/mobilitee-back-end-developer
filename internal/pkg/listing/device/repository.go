package device

import (
	"github.com/jinzhu/gorm"
)

// Repository is the handles the persistency layer
// for listing Devices.
type Repository struct {
	db *gorm.DB
}

// NewRepository returns the persistency layer for
// listing Devices.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}
