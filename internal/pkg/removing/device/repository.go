package device

import "github.com/jinzhu/gorm"

// Repository implements the Remover interface
// for removeing devices from persistency layer
type Repository struct {
	db *gorm.DB
}

// NewRepository returns a reference to Repository
// for removing devices
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

// Remove deletes the device from persistency layer
func (r *Repository) Remove(d *Device) error {
	return r.db.Delete(d).Error
}
