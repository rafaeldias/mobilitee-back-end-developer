package user

import "time"

// DeviceReader reads informations about the user's devices.
type DeviceReader interface {
	LatestRemoved(user int) (*Device, error)
	LatestExchange(user int) (*Device, error)
	Count(user int) (int, error)
}

// Usecase controls the flow for reading data
// about user exchanging statuses.
type Usecase struct {
	repo DeviceReader
}

// New returns a reference to Usecase.
func New(repo DeviceReader) *Usecase {
	return &Usecase{repo}
}

//IsExchanging returns if user is trying to exchange a device and its id or error.
func (u *Usecase) IsExchanging(user int) (exchanging bool, device int, err error) {
	var d *Device

	if d, err = u.repo.LatestRemoved(user); err != nil {
		return
	}

	if exchanging = d.BeingExchanged(); exchanging {
		device = d.ID
	}

	return
}

// LatestExchangeExpiresAt returns the next date the user can do a new exchanging
// if the lastest exchanging hasn't expired or never occurred, or error.
func (u *Usecase) LatestExchangeExpiresAt(user int) (time.Time, error) {
	d, err := u.repo.LatestExchange(user)
	if err != nil {
		return time.Time{}, err
	}

	return d.LatestExchangeExpiresAt(), nil
}

// CountDevices returns the number of devices a user created
func (u *Usecase) CountDevices(user int) (int, error) {
	return u.repo.Count(user)
}
