package device

import (
	"reflect"
	"testing"
)

type repoWriter struct {
	id    int
	err   error
	dvice *Device
}

func (r *repoWriter) Write(d *Device) (id int, err error) {
	id = r.id
	err = r.err

	return
}

func TestNew(t *testing.T) {
	testCases := []struct {
		repo Writer
	}{
		{nil},
		{&repoWriter{}},
	}

	for _, tc := range testCases {
		d := New(tc.repo)

		if !reflect.DeepEqual(d.repo, tc.repo) {
			t.Errorf("got New(%+v): %+v, want: %+v", tc.repo, d.repo, tc.repo)
		}
	}
}

func TestWriter(t *testing.T) {
	testCases := []struct {
		id    int
		dvice *Device
	}{
		{1, &Device{Name: "Test", Model: "Android", User: 1}},
	}

	for _, tc := range testCases {
		repo := &repoWriter{id: tc.id}

		device := New(repo)

		id, err := device.Write(tc.dvice)
		if err != nil {
			t.Errorf("got error while calling Device.Writer(%+v): %s, want nil",
				tc.dvice, err.Error())
		}

		if tc.id != id {
			t.Errorf("got devices from Device.Write(%+v): %d, want: %d", tc.dvice,
				id, tc.id)
		}
	}
}

func TestValidFieldsError(t *testing.T) {
	testCases := []struct {
		device *Device
	}{
		{&Device{Name: "", User: 1}},
		{&Device{Name: "Testing", User: 0}},
	}

	for _, tc := range testCases {
		u := New(&repoWriter{})
		err := u.ValidFields(tc.device)
		if err == nil {
			t.Errorf("got error nil while calling Usecase.ValidFields(%+v), want not nil",
				tc.device)
		}
	}
}
