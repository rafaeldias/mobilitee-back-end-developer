package device

import (
	"testing"
	//"time"
)

func TestDeviceValidFieldsError(t *testing.T) {
	testCases := []struct {
		Device         *Device
	}{
		{&Device{Name: "", Model: "Android"}},
		{&Device{Name: "Testing", Model: ""}},
		{&Device{Name: "Testing", Model: "Windows"}},
	}

	for _, tc := range testCases {
		err := tc.Device.ValidFields()
		if err == nil {
			t.Errorf("got error nil while calling Device.ValidFields(%+v), want not nil",
				tc.Device)
		}
	}
}

//func TestDeviceValidError(t *testing.T) {
//	testCases := []struct {
//		Device         *Device
//		Devices        []*Device
//		LastExchangeAt time.Time
//		LastRemovedAt  time.Time
//	}{
//		{
//			&Device{Name: "Testing", Model: "Android"},
//			[]*Device{
//				&Device{Name: "Smartphone", Model: "Android"},
//				&Device{Name: "Smartphone", Model: "iOS"},
//				&Device{Name: "Tablet", Model: "iOS"},
//			},
//			time.Time{},
//			time.Time{},
//		},
//		{
//			&Device{Name: "Testing", Model: "Android"},
//			[]*Device{
//				&Device{Name: "Smartphone", Model: "Android"},
//				&Device{Name: "Smartphone", Model: "iOS"},
//				&Device{Name: "Tablet", Model: "iOS"},
//			},
//			time.Now().AddDate(0, 0, -20),
//		},
//		//{
//		//	&Device{Name: "Testing", Model: "Android"},
//		//	[]*Device{
//		//		&Device{Name: "Smartphone", Model: "Android"},
//		//		&Device{Name: "Smartphone", Model: "iOS"},
//		//	},
//		//	time.Now().AddDate(0, 0, -10),
//		//},
//	}
//
//	for _, tc := range testCases {
//		err := tc.Device.Valid(tc.Devices, tc.LastExchange)
//		fmt.Println(err.Error())
//		if err == nil {
//			t.Errorf("got error nil while calling Device.Valid(%+v, %+v), want not nil",
//				tc.Devices, tc.LastExchange)
//		}
//	}
//}
