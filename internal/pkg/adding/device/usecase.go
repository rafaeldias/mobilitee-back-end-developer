package device

import "time"

// Writer is the interface that writes the Device.
type Writer interface {
	// returns the id the new Device or error.
	Write(d *Device) (id int, err error)
}

// WriteExchanger writes a device or exchanges it for another.
type WriteExchanger interface {
	Writer
	// returns the id of the the new device or error
	Exchange(device int, nw *Device) (id int, err error)
}

// User is the interface that reads user data statuses.
type User interface {
	IsExchanging(user int) (exchanging bool, device int, err error)
	LatestExchangeExpiresAt(user int) (time.Time, error)
	CountDevices(user int) (int, error)
}

// Usecase controls the data flow for creating a device.
type Usecase struct {
	repo WriteExchanger
	user User
}

// New returns a pointer to Usecase.
func New(repo WriteExchanger, user User) *Usecase {
	return &Usecase{repo, user}
}

// Write controls the flow of data for creating/exchanging a device.
func (u *Usecase) Write(dvice *Device) (id int, err error) {
	var (
		exchanging		bool
		old			int
		latestExchangeExpiresAt time.Time
		userDevices		int
	)

	if err = u.validFields(dvice); err != nil {
		return
	}

	if exchanging, old, err = u.user.IsExchanging(dvice.User); err != nil {
		return
	}

	if latestExchangeExpiresAt, err = u.user.LatestExchangeExpiresAt(dvice.User); err != nil {
		return
	}

	if userDevices, err = u.user.CountDevices(dvice.User); err != nil {
		return
	}

	if err = dvice.Valid(userDevices, exchanging, latestExchangeExpiresAt); err != nil {
		return
	}

	if exchanging {
		id, err = u.repo.Exchange(old, dvice)
		return
	}

	id, err = u.repo.Write(dvice)

	return
}

// validFields validates application-specific requirements
// for creating/exchanging a device.
func (u *Usecase) validFields(d *Device) error {
	if d.Name == "" {
		return &InvalidError{"attribute `Name` must not be empty"}
	}

	if d.User == 0 {
		return &InvalidError{"invalid user"}
	}

	return nil
}
