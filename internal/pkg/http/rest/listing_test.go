package rest

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	//"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/listing/device"
)

type getterTest struct {
	paths  []string
	handle httprouter.Handle
}

func (g *getterTest) GET(path string, h httprouter.Handle) {
	g.paths = append(g.paths, path)
	g.handle = h
}

// Don't need this yer
//type deviceReader struct {
//	id      int
//	err     error
//	devices []*device.Device
//}
//
//func (r *deviceReader) Read(id int) ([]*device.Device, error) {
//	return r.devices, r.err
//
//}

func TestGetDeviceRoutes(t *testing.T) {
	testCases := []struct {
		path string
	}{
		{"/devices"},
		{"/devices/:id"},
	}

	for _, tc := range testCases {
		getter := &getterTest{}
		GetDevices(getter, nil)

		var pathInList bool

		for _, path := range getter.paths {
			if pathInList = path == tc.path; pathInList {
				break
			}
		}

		if !pathInList {
			t.Errorf("got path `%s` out of list, want it in", tc.path)
		}
	}
}
