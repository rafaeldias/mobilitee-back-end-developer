package device

import (
	"fmt"
	"testing"
	"time"
)

func TestDeviceValidError(t *testing.T) {
	testCases := []struct {
		Device       *Device
		Devices      []*Device
		LastExchange time.Time
	}{
		{
			&Device{Name: "Testing", Model: "Android"},
			[]*Device{
				&Device{Name: "Smartphone", Model: "Android"},
				&Device{Name: "Smartphone", Model: "iOS"},
				&Device{Name: "Tablet", Model: "iOS"},
			},
			time.Time{},
		},
		{
			&Device{Name: "Testing", Model: "Android"},
			[]*Device{
				&Device{Name: "Smartphone", Model: "Android"},
				&Device{Name: "Smartphone", Model: "iOS"},
				&Device{Name: "Tablet", Model: "iOS"},
			},
			time.Now().AddDate(0, 0, -20),
		},
		//{
		//	&Device{Name: "Testing", Model: "Android"},
		//	[]*Device{
		//		&Device{Name: "Smartphone", Model: "Android"},
		//		&Device{Name: "Smartphone", Model: "iOS"},
		//	},
		//	time.Now().AddDate(0, 0, -10),
		//},
	}

	for _, tc := range testCases {
		err := tc.Device.Valid(tc.Devices, tc.LastExchange)
		if err == nil {
			t.Errorf("got error nil while calling Device.Valid(%+v, %+v), want not nil",
				tc.Devices, tc.LastExchange)
		}
	}
}
