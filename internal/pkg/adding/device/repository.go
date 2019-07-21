package device

import "github.com/jinzhu/gorm"

// Repository implements the WriteExchanger interface
// for writing/exchanging devices in persistency layer
type Repository struct {
	db *gorm.DB
}

// NewRepository returns a reference to Repository
// for adding/exchanging devices
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

// Write saves the Device in the persistenct layer
func (r *Repository) Write(d *Device) (id int, err error) {
	err = r.db.Omit("exchanged").Create(d).Error
	id = d.ID

	return
}

// Exchange writes a new device and updates de old one for flagging it
// with the id of the newer created one.
func (r *Repository) Exchange(old *Device, nw *Device) (id int, err error ) {
	tx := r.db.Begin()

	if err = tx.Error; err != nil {
		return
	}

	if err = tx.Omit("exchanged").Create(nw).Error; err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Model(old).Update("exchanged", nw.ID).Error; err != nil {
		tx.Rollback()
		return
	}

	id = nw.ID

	err = tx.Commit().Error

	return
}
