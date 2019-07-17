package device

import (
	"reflect"
	"testing"
)

type repoUpdater struct {
	err    error
	ID     int
	Device *Device
}

func (u *repoUpdater) Update(ID int, d *Device) error {
	u.ID = ID
	u.Device = d
	return u.err
}

func TestNew(t *testing.T) {
	testCases := []struct {
		repo Updater
	}{
		{nil},
		{&repoUpdater{}},
	}

	for _, tc := range testCases {
		d := New(tc.repo)

		if !reflect.DeepEqual(d.repo, tc.repo) {
			t.Errorf("got New(%+v): %+v, want: %+v", tc.repo, d.repo, tc.repo)
		}
	}
}

func TestDeviceUpdate(t *testing.T) {
	testCases := []struct {
		ID     int
		Device *Device
		repo   *repoUpdater
	}{
		{
			1,
			&Device{"New Name Testing"},
			&repoUpdater{},
		},
	}

	for _, tc := range testCases {
		d := New(tc.repo)

		if err := d.Update(tc.ID, tc.Device); err != nil {
			t.Errorf("go error while calling device.Update(%d, %+v): %s, want nil",
				tc.ID, tc.Device, err.Error())
		}

		if !reflect.DeepEqual(tc.repo.Device, tc.Device) {
			t.Errorf("got repo.Device: %+v, want %+v", tc.repo.Device, tc.Device)
		}

		if tc.repo.ID != tc.ID {
			t.Errorf("got repo.ID: %d, want %d", tc.repo.ID, tc.ID)
		}
	}
}

func TestDeviceUpdateError(t *testing.T) {
	testCases := []struct {
		ID     int
		Device *Device
		repo   *repoUpdater
	}{
		{
			0,
			&Device{"New Name Testing"},
			&repoUpdater{},
		},
		{
			1,
			&Device{},
			&repoUpdater{},
		},
	}

	for _, tc := range testCases {
		d := New(tc.repo)

		if err := d.Update(tc.ID, tc.Device); err == nil {
			t.Errorf("go error nil while calling device.Update(%d, %+v), want not nil",
				tc.ID, tc.Device)
		}
	}
}
