package user

import (
	"github.com/jinzhu/gorm"
)

// Repository is the persistency layer for reading
// information of User's operations.
type Repository struct {
	db *gorm.DB
}

// NewRepository returns a reference to the Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

// Returns a reference to the latest removed device by user.
func (r *Repository) LatestRemoved(user int) (*Device, error) {
	d := &Device{}

	if err := r.db.Select([]string{"id", "MAX(deleted_at) as latest_removed_at"}).
		Where(`"user" = ? AND exchanged IS NULL AND deleted_at IS NOT NULL`, user).
		Group("id").
		Order("deleted_at ASC").
		Limit(1).
		Find(d).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return d, nil
}

// LatestExchange returns the latest exchanged device by user or error.
func (r *Repository) LatestExchange(user int) (*Device, error) {
	d := &Device{}

	if err := r.db.Select("MAX(devices.created_at) as latest_exchange_at").
		Joins("JOIN devices d ON d.exchanged = devices.id").
		Where("devices.user = ?", user).
		Find(d).Error; err != nil && !gorm.IsRecordNotFoundError(err)  {
		return nil, err
	}
	return d, nil
}

// Count returns the number of user's active devices.
func (r *Repository) Count(user int) (int, error) {
	var total int
	if err := r.db.Model(&Device{}).
		Select("count(*)").
		Where(`"user" = ? AND deleted_at IS NULL`, user).
		Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

