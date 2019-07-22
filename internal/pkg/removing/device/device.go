package device

import (
	"time"
	"fmt"
)

// Error implements error interface
type InvalidOperation struct {
	Msg string
}

func (e *InvalidOperation) Error() string {
	return e.Msg
}

// Device is a representation of a device that will
// be removed
type Device struct {
	ID int
	User int
	DeletedAt time.Time
}

func (d *Device) Valid(devices int, latestExchangeExpiresAt time.Time) error {
	if devices == 1 && !latestExchangeExpiresAt.IsZero() {
		return &InvalidOperation{fmt.Sprintf(
			"You can't remove the last device because you've recently made a device exchange. You latest exchange expires at %s",
			latestExchangeExpiresAt.Format("2006-01-02"))}
	}
	return nil
}
