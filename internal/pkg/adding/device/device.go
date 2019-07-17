package device

import (
	"errors"
	"fmt"
	"time"
)

// MaxDevicesExchangeNotExpired occurs when the limit
// of devices has been reached and an exchange has not
// expired yet.
//type MaxDevicesExchangeNotExpired struct {
//	expiresIn time.Duration
//}
//
//func (e *MaxDevicesExchangeNotExpired) Error() {
//	fmt.Sprintf(
//		"devices limit exceeded and you cannot do an exchanging. You can exchange a device in %d days",
//		e.expiresIn,
//	)
//}

// Device is the entity that will validate the business
// rules for creating a user device in persistency layer
type Device struct {
	Name      string
	Model     string
	User      int
	CreatedAt time.Time
}

// Valid returns error if the current device isn't able
// to be written to the persistency layer
func (d *Device) Valid(ds []*Device, lastExchange time.Time) error {
	exchangeExpired := true

	if !lastExchange.IsZero() {
		exchangeExpired = time.Since(lastExchange)/(24*time.Hour) > 30
	}

	if len(ds) == 3 {
		if exchangeExpired {
			return errors.New("devices max limit exceeded, but you still can do an exchanging")
		}

		return fmt.Errorf(
			"devices max limit exceeded and you cannot do an exchanging. You can exchange a device in %d days",
			(time.Since(lastExchange) / (24 * time.Hour)),
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
