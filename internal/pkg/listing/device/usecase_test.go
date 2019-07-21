package device

import (
	"errors"
	"reflect"
	"testing"
)

type repoReader struct {
	err error
}

func (r *repoReader) Read(id int) ([]*Device, error) {
	return []*Device{&Device{ID: id}}, r.err
}

func TestNew(t *testing.T) {
	testCases := []struct {
		repo Reader
	}{
		{nil},
		{&repoReader{}},
	}

	for _, tc := range testCases {
		d := New(tc.repo)

		if !reflect.DeepEqual(d.repo, tc.repo) {
			t.Errorf("got New(%+v): %+v, want: %+v", tc.repo, d.repo, tc.repo)
		}
	}
}

func TestUsecaseRead(t *testing.T) {
	testCases := []struct {
		id   int
		want []*Device
	}{
		{1, []*Device{&Device{ID: 1}}},
		{2, []*Device{&Device{ID: 2}}},
	}

	for _, tc := range testCases {
		repo := &repoReader{}

		device := New(repo)

		devices, err := device.Read(tc.id)
		if err != nil {
			t.Errorf("got error while calling Device.Read(%d): %s, want nil",
				tc.id, err.Error())
		}

		if !reflect.DeepEqual(tc.want, devices) {
			t.Errorf("got devices from Device.Read(%d): %+v, want: %+v", tc.id,
				devices, tc.want)
		}
	}
}

func TestUsecaseReadError(t *testing.T) {
	testCases := []struct {
		repo Reader
	}{
		{&repoReader{errors.New("Testing Usecase Error")}},
	}

	for _, tc := range testCases {
		device := New(tc.repo)

		_, err := device.Read(0)
		if err == nil {
			t.Error("got error nil while calling Device.Read(0):, want not nil")
		}
	}
}
