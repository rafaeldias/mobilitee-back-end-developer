package user

import (
	"testing"
	"time"
)

func TestBeingExchanged(t *testing.T) {
	testCases := []struct{
		dvice *Device
		want bool
	}{
		{
			&Device{LatestRemovedAt: time.Now().AddDate(0, 0, -20)},
			true,
		},
		{
			&Device{LatestRemovedAt: time.Now().AddDate(0, 0, -31)},
			false,
		},
		{
			&Device{LatestRemovedAt: time.Time{}},
			false,
		},
	}

	for _, tc := range testCases {
		if got := tc.dvice.BeingExchanged(); got != tc.want {
			t.Errorf("got %t when calling Device.IsExchanging(), want %t", got, tc.want)
		}
	}
}

func TestLatestExchangeExpiresAt(t *testing.T) {
	testCases := []struct{
		dvice *Device
		want time.Time
	}{
		{
			&Device{LatestExchangeAt: time.Now().AddDate(0, 0, -10)},
			time.Now().Add(time.Hour * 24 * 30 - time.Since(time.Now().AddDate(0, 0, -10))),
		},
		{
			&Device{LatestExchangeAt: time.Now().AddDate(0, 0, -31)},
			time.Time{},
		},
		{
			&Device{LatestExchangeAt: time.Time{}},
			time.Time{},
		},
	}

	for _, tc := range testCases {
		if got := tc.dvice.LatestExchangeExpiresAt(); got != tc.want {
			t.Errorf("got %s when calling Device.LatestExchangeExpiresAt(), want %s", got, tc.want)
		}
	}
}
