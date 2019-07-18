package device

import (
	"reflect"
	"testing"
)

type repoUpdater struct {
	err    error
	id     int
	device *Device
}

func (u *repoUpdater) Update(id int, d *Device) error {
	u.id = id
	u.device = d
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
		id     int
		device *Device
		repo   *repoUpdater
	}{
		{
			1,
			&Device{Name: "New Name Testing"},
			&repoUpdater{},
		},
	}

	for _, tc := range testCases {
		d := New(tc.repo)

		if err := d.Update(tc.id, tc.device); err != nil {
			t.Errorf("go error while calling device.Update(%d, %+v): %s, want nil",
				tc.id, tc.device, err.Error())
		}

		if !reflect.DeepEqual(tc.repo.device, tc.device) {
			t.Errorf("got repo.Device: %+v, want %+v", tc.repo.device, tc.device)
		}

		if tc.repo.id != tc.id {
			t.Errorf("got repo.ID: %d, want %d", tc.repo.id, tc.id)
		}
	}
}

func TestDeviceUpdateError(t *testing.T) {
	testCases := []struct {
		id     int
		device *Device
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

		if err := d.Update(tc.id, tc.device); err == nil {
			t.Errorf("go error nil while calling device.Update(%d, %+v), want not nil",
				tc.id, tc.device)
		}
	}
}
