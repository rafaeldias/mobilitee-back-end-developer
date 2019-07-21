package device

import (
	"testing"
	"time"
)

func TestValid(t *testing.T) {
	testCases := []struct {
		device            *Device
		devices           int
		beingExchanged    bool
		lastestExchangeAt time.Time
	}{
		{
			&Device{Model: "Android"},
			0,
			false,
			time.Time{},
		},
		{
			&Device{Model: "iOS"},
			2,
			true,
			time.Time{},
		},
	}

	for _, tc := range testCases {
		err := tc.device.Valid(tc.devices, tc.beingExchanged, tc.lastestExchangeAt)
		if err != nil {
			t.Errorf("got error while calling Device.Valid(%d, %t, %s): %s, want nil",
				tc.devices, tc.beingExchanged, tc.lastestExchangeAt, err.Error())
		}
	}
}

func TestValidError(t *testing.T) {
	testCases := []struct {
		device            *Device
		devices           int
		beingExchanged    bool
		lastestExchangeAt time.Time
	}{
		{
			&Device{Model: "windows"},
			0,
			false,
			time.Time{},
		},
		{
			&Device{Model: "Android"},
			3,
			false,
			time.Time{},
		},
		{
			&Device{Model: "iOS"},
			3,
			false,
			time.Now().AddDate(0, 0, 20),
		},
		{
			&Device{Model: "Android"},
			2,
			true,
			time.Now().AddDate(0, 0, 20),
		},
	}

	for _, tc := range testCases {
		err := tc.device.Valid(tc.devices, tc.beingExchanged, tc.lastestExchangeAt)
		if err == nil {
			t.Errorf("got error nil while calling Device.Valid(%d, %t, %s), want not nil",
				tc.devices, tc.beingExchanged, tc.lastestExchangeAt)
		}
	}
}
