package device

import "time"

// Remover is the interface that removes
// the device from pesistency layer
type Remover interface {
	Remove(d *Device) error
}

// User is the interface that reads user data statuses.
type User interface {
	LatestExchangeExpiresAt(user int) (time.Time, error)
	CountDevices(user int) (int, error)
}

// Usecase controls the flow of data for removing a device
type Usecase struct {
	repo Remover
	user User
}

// New returns a reference to the Usecase
func New(repo Remover, user User) *Usecase {
	return &Usecase{repo, user}
}

// Remove validates and removes the device
func (u *Usecase) Remove(d *Device) error {
	latestExchangeExpiresAt, err := u.user.LatestExchangeExpiresAt(d.User)
	if err != nil {
		return err
	}
	userDevices, err := u.user.CountDevices(d.User);
	if err != nil {
		return err
	}

	if err := d.Valid(userDevices, latestExchangeExpiresAt); err != nil {
		return err
	}

	return u.repo.Remove(d)
}
