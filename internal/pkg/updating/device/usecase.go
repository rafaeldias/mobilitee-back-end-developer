package device

import "fmt"

// ErrInvalidInput describes the invalid data user inputed
type ErrInvalidInput struct {
	Attr string
	Must string
}

func (e *ErrInvalidInput) Error() string {
	return fmt.Sprintf("attibute %s must %s", e.Attr, e.Must)
}

// ErrNotFound describes a device being edited not found
type ErrNotFound struct {
	ID int
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("device of ID %d was not found", e.ID)
}

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
		return &ErrInvalidInput{Attr: "ID", Must: "be greater than 0"}
	}

	if d.Name == "" {
		return &ErrInvalidInput{Attr: "Name", Must: "not be empty"}
	}

	return u.repo.Update(id, d)
}
