package device

import (
	"testing"
	"time"
)

func TestDeviceValid(t *testing.T) {
	testCases := []struct {
		devices int
		latestExchangeExperesAt time.Time
	}{
		{1, time.Time{}},
		{2, time.Time{}},
		{2, time.Now().AddDate(0, 0, -10)},
	}

	for _, tc := range testCases {
		d := &Device{}
		if err := d.Valid(tc.devices, tc.latestExchangeExperesAt); err != nil {
			t.Fatalf("got error while calling Device.Valid(%d, %s): %s; want nil",
				tc.devices, tc.latestExchangeExperesAt, err.Error())
		}
	}
}

func TestDeviceValidError(t *testing.T) {
	testCases := []struct {
		devices int
		latestExchangeExperesAt time.Time
	}{
		{1, time.Now().AddDate(0, 0, -20)},
	}

	for _, tc := range testCases {
		d := &Device{}
		if err := d.Valid(tc.devices, tc.latestExchangeExperesAt); err == nil {
			t.Fatalf("got error nil while calling Device.Valid(%d, %s); want not nil",
				tc.devices, tc.latestExchangeExperesAt)
		}
	}
}
