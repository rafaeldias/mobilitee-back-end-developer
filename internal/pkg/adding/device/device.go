package device

import (
	"fmt"
	"time"
)

// InvalidError describes any error occurred by invalid input
type InvalidError struct {
	msg string
}

// Implements the error interface
func (e *InvalidError) Error() string {
	return e.msg
}

// Device is the entity that will validate the business
// rules for creating a user device in persistency layer
type Device struct {
	ID	  int
	Name      string
	Model     string
	User      int
	Exchanged int
	CreatedAt time.Time
}


// Valid returns error if the current device is able to be written/exchanged
// to the persistency layer based on business logic
func (d *Device) Valid(devices int, beingExchanged bool, latestExchangeExpiresAt time.Time) error {
	if err := d.validModel(); err != nil {
		return err
	}

	latestExchangeExpired := latestExchangeExpiresAt.IsZero()

	if devices == 3 {
		if latestExchangeExpired {
			return &InvalidError{"devices max limit exceeded, but you still can do an exchanging."}
		}

		return &InvalidError{fmt.Sprintf(
			"devices max limit exceeded and an exchange cannot be made. A new exchange can be made as of %s.",
			latestExchangeExpiresAt.Format("2006-01-02"))}
	}

	if beingExchanged && !latestExchangeExpired {
		return &InvalidError{fmt.Sprintf(
			"you've recently made a device exchange. You can make another one as of %s",
			latestExchangeExpiresAt.Format("2006-01-02"))}
	}

	return nil
}

func (d *Device) validModel() error {
	if d.Model != "Android" && d.Model != "iOS" {
		return &InvalidError{"attribute `Model` must be `Android` or `iOS`"}
	}

	return nil
}
