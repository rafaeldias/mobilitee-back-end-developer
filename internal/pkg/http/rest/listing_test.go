package rest

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/listing/device"
)

type getterTest struct {
	paths  []string
	handle httprouter.Handle
}

func (g *getterTest) GET(path string, h httprouter.Handle) {
	g.paths = append(g.paths, path)
	g.handle = h
}

type deviceReader struct {
	id      int
	err     error
	devices []*device.Device
}

func (r *deviceReader) Read(id int) ([]*device.Device, error) {
	return r.devices, r.err

}

func TestGetDeviceRoutes(t *testing.T) {
	testCases := []struct {
		path   string
		getter *getterTest
		reader *deviceReader
	}{
		{
			"/devices",
			&getterTest{},
			&deviceReader{},
		},
		{
			"/devices/:id",
			&getterTest{},
			&deviceReader{},
		},
	}

	for _, tc := range testCases {
		GetDevices(tc.getter, tc.reader)

		var pathInList bool

		for _, path := range tc.getter.paths {
			if pathInList = path == tc.path; pathInList {
				break
			}
		}

		if !pathInList {
			t.Errorf("got path `%s` out of list, want it in", tc.path)
		}
	}
}
