package device

import "errors"

// ErrInvalidInput describes the invalid data user inputed
type ErrInvalidInput error

// ErrDeviceNotFound describes a device being edited not found
type ErrDeviceNotFound error

// Updater edit the Device
type Updater interface {
	Update(id int, d *Device) error
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
func (u *Usecase) Update(id int, d *Device) error {
	if id == 0 {
		return ErrInvalidInput(errors.New("ID attribute must be greater than 0"))
	}

	if d.Name == "" {
		return ErrInvalidInput(errors.New("Name attribute must not be empty"))
	}

	return u.repo.Update(id, d)
}
