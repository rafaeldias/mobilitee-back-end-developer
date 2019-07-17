package device

import (
	"errors"
	"fmt"
	"time"
)

// InvalidError describes any error occurred by invalid input
type InvalidError struct {
	msg string
}

func (e *MaxDevicesExchangeNotExpired) Error() {
	return e.msg
}

// Device is the entity that will validate the business
// rules for creating a user device in persistency layer
type Device struct {
	Name      string
	Model     string
	User      int
	CreatedAt time.Time
}

// ValidFields validates the fields of the Device
func (d *Device) ValidFields() error {
		if d.Name == "" {
			return &InvalidError{"attribute `Name` must not be empty"}
		}

		if d.Model == "" {
			return &InvalidError{"attribute `Model` must not be empty"}
		}

		if d.Model != "Android" && d.Model != "iOS" {
			return &InvalidError{"attribute `Model` must be `Android` or `iOS`"}
		}

		return nil
}

// Valid returns error if the current device isn't able
// to be written to the persistency layer
func (d *Device) Valid(ds []*Device, lastExchangedAt, lastRemovedAt time.Time) error {
	exchangeExpired := true

	if !lastExchangedAt.IsZero() {
		exchangeExpired = time.Since(lastExchangedAt)/(24*time.Hour) > 30
	}

	if len(ds) == 3 {
		if exchangeExpired {
			return &InvalidError{"devices max limit exceeded, but you still can do an exchanging"}
		}

		return &InvalidError{fmt.Errorf(
			"devices max limit exceeded and you cannot do an exchanging. You can exchange a device at %s",
			(time.Since(lastExchangedAt) / (24 * time.Hour))},
		)
	}

	//if !exchangeExpired {
	//	return fmt.Errorf(
	//		"You've recently exchanged a device. You can exchange or add a new one in %d days",
	//		(time.Since(lastExchange) / (24 * time.Hour)),
	//	)
	//}

	return nil
}
