package device

import "github.com/jinzhu/gorm"

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

// Read returns a list of devices or error from the
// persistency layer
func (r *Repository) Read(id int) ([]*Device, error) {
	var args []interface{}

	if id > 0 {
		args = append(args, id)
	}

	devices := []*Device{}

	err := r.db.Find(&devices, args...).Error

	return devices, err
}
