package device

import (
	"fmt"
	"testing"
	"time"
)

func TestDeviceValidFieldsError(t *testing.T) {
	testCases := []struct {
		Device *Device
	}{
		{&Device{Name: "", Model: "Android", User: 1}},
		{&Device{Name: "Testing", Model: "", User: 2}},
		{&Device{Name: "Testing", Model: "Windows", User: 3}},
		{&Device{Name: "Testing", Model: "iOS"}},
	}

	for _, tc := range testCases {
		err := tc.Device.ValidFields()
		if err == nil {
			t.Errorf("got error nil while calling Device.ValidFields(%+v), want not nil",
				tc.Device)
		}
	}
}

func TestDeviceValidError(t *testing.T) {
	testCases := []struct {
		Device          *Device
		Devices         []*Device
		LastExchangedAt time.Time
		LastRemovedAt   time.Time
	}{
		{
			&Device{Name: "Testing", Model: "Android"},
			[]*Device{
				&Device{Name: "Smartphone", Model: "Android"},
				&Device{Name: "Smartphone", Model: "iOS"},
				&Device{Name: "Tablet", Model: "iOS"},
			},
			time.Time{},
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
			time.Time{},
		},
		//{
		//	&Device{Name: "Testing", Model: "Android"},
		//	[]*Device{
		//		&Device{Name: "Smartphone", Model: "Android"},
		//		&Device{Name: "Smartphone", Model: "iOS"},
		//	},
		//	time.Now().AddDate(0, 0, -10),
		//	time.Now().AddDate(0, 0, -5),
		//},
	}

	for _, tc := range testCases {
		err := tc.Device.Valid(tc.Devices, tc.LastExchangedAt, tc.LastRemovedAt)
		fmt.Println(err.Error())
		if err == nil {
			t.Errorf("got error nil while calling Device.Valid(%+v, %+v), want not nil",
				tc.Devices, tc.LastExchangedAt)
		}
	}
}
