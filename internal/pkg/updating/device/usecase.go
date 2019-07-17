package device

import "errors"

// Updater edit the Device
type Updater interface {
	Update(ID int, d *Device) error
}

// Usecase controls the flow of data for updating the Device
type Usecase struct {
	// Repository should implement the Updater interface
	repo Updater
}

// New returns a new usecase for updating devices
func New(repo Updater) *Usecase {
	return &Usecase{repo}
}

// Update validates and edits a Device through the persistency layer
func (u *Usecase) Update(ID int, d *Device) error {
	if ID == 0 {
		return errors.New("ID attribute must be greater than 0")
	}

	if d.Name == "" {
		return errors.New("Name attribute must not be empty")
	}

	return u.repo.Update(ID, d)
}
