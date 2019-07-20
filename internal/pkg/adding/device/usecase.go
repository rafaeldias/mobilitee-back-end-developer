package device

import "time"

// Writer is the interface that writes the Device in
// Persistency layer
type Writer interface {
	// Returns the id the new Device
	Write(d *Device) (id int, err error)
}

// Exchanger is
type Exchanger interface {
	IsExchanging(user int) (bool, error)
	LastExchangeExpired(user int) (expired bool, nexExchangeAt time.Time, err error)
	CountDevices(user int) int
}

// Usecase controls the data flow for creating a device
type Usecase struct {
	repo Writer
}

// New returns a pointer to Usecase
func New(repo Writer) *Usecase {
	return &Usecase{repo}
}

// Write is
func (u *Usecase) Write(d *Device) (id int, err error) {
	id, err = u.repo.Write(d)

	return
}

// ValidFields validates application-specific requirements
// for creating a device
func (u *Usecase) ValidFields(d *Device) error {
	if d.Name == "" {
		return &InvalidError{"attribute `Name` must not be empty"}
	}

	if d.User == 0 {
		return &InvalidError{"invalid user"}
	}

	return nil
}
