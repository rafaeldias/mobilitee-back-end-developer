package device

import "github.com/jinzhu/gorm"

// Repository implements the Updater interface
// for editing devices in persistency layer
type Repository struct {
	db *gorm.DB
}

// NewRepository returns a reference to Device Updating Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

// Update impelements the Updater interface
func (r *Repository) Update(id int, d *Device) error {
	err := r.db.Model(Device{}).Where("id = ?", id).Update(d).Error
	if gorm.IsRecordNotFoundError(err) {
		return &ErrNotFound{id}
	}
	return err
}
