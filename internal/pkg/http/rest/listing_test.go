package rest

import (
	"testing"
	"net/http/httptest"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/listing/device"
)

type getterTest struct {
	paths  []string
	handle map[string]httprouter.Handle
}

func (g *getterTest) GET(path string, h httprouter.Handle) {
	g.paths = append(g.paths, path)
	if g.handle == nil {
		g.handle = map[string]httprouter.Handle{}
	}
	g.handle[path] = h
}


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

type deviceReader struct {
	id      int
	err     error
	devices []*device.Device
}

func (r *deviceReader) Read(id int) ([]*device.Device, error) {
	r.id = id
	return r.devices, r.err

}

func TestGetDevicesHandle(t *testing.T) {
	testCases := []struct {
		path           string
		params         httprouter.Params
		getter         *getterTest
		reader         *deviceReader
		wantStatusCode int
	}{
		{
			"/devices",
			httprouter.Params{},
			&getterTest{},
			&deviceReader{},
			http.StatusOK,
		},
	}

	for _, tc := range testCases {
		GetDevices(tc.getter, tc.reader)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)

		tc.getter.handle[tc.path](rw, req, tc.params)

		r := rw.Result()

		if r.StatusCode != tc.wantStatusCode {
			t.Errorf("got http status code %d, want %d", r.StatusCode, tc.wantStatusCode)
		}
	}
