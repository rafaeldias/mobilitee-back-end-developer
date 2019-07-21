package user

import (
	"time"
)

// Device is the representation of a user Device.
type Device struct {
	ID		  int
	Total		  int
	LatestExchangeAt  time.Time
	LatestRemovedAt   time.Time
}

// BeingExchanged returns if is being exchanged by the user based on the
// time difference from the last device he/she has removed.
func (d *Device) BeingExchanged() bool {
	if (time.Since(d.LatestRemovedAt).Hours() / 24) < 30 {
		return true
	}
	return false
}

// LatestExchangeExpiresAt returns the date the user will be able to to another exchange.
// If his/her the latest exchange expired or an exchange never occured yet, a zero time is returned.
func (d *Device) LatestExchangeExpiresAt() time.Time {
	if deltaLatestExchange := time.Since(d.LatestExchangeAt); (deltaLatestExchange.Hours() / 24) < 30 {
		thirtyDays := time.Hour * 24 * 30
		return time.Now().Add(thirtyDays - deltaLatestExchange)
	}

	return time.Time{}
}
