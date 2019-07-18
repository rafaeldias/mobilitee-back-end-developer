package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

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

type deviceReader struct {
	id      int
	err     error
	devices []*device.Device
}

func (r *deviceReader) Read(id int) ([]*device.Device, error) {
	r.id = id
	return r.devices, r.err
}

func TestGetDevices(t *testing.T) {
	testCases := []struct {
		path            string
		params          httprouter.Params
		getter          *getterTest
		reader          *deviceReader
		wantDevices     *Devices
		wantDevice      *device.Device
		wantContentType string
		wantStatusCode  int
	}{
		{
			"/devices",
			httprouter.Params{},
			&getterTest{},
			&deviceReader{
				devices: []*device.Device{
					{ID: 1, Name: "Testing"},
				},
			},
			&Devices{
				Data: []*device.Device{
					{ID: 1, Name: "Testing"},
				},
			},
			nil,
			"application/json",
			http.StatusOK,
		},
		{
			"/devices",
			httprouter.Params{},
			&getterTest{},
			&deviceReader{
				devices: []*device.Device{},
			},
			&Devices{
				Data: []*device.Device{},
			},
			nil,
			"application/json",
			http.StatusOK,
		},
		{
			"/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			&getterTest{},
			&deviceReader{
				devices: []*device.Device{
					{ID: 1, Name: "Testing"},
				},
			},
			nil,
			&device.Device{ID: 1, Name: "Testing"},
			"application/json",
			http.StatusOK,
		},
	}

	for _, tc := range testCases {
		GetDevices(tc.getter, tc.reader)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)

		tc.getter.handle[tc.path](rw, req, tc.params)

		res := rw.Result()

		if res.StatusCode != tc.wantStatusCode {
			t.Errorf("got http status code %d, want %d", res.StatusCode, tc.wantStatusCode)
		}

		if c := rw.Header().Get("Content-Type"); c != tc.wantContentType {
			t.Errorf("got http header `Content-Type`: %s, want %s", c, tc.wantContentType)
		}

		// Test for GET /devices response
		if tc.wantDevices != nil {

			if tc.reader.id != 0 {
				t.Errorf("got Read() param: %d, want 0", tc.reader.id)
			}

			var devices *Devices

			if err := json.NewDecoder(res.Body).Decode(&devices); err != nil {
				t.Errorf("got error while decoding response body: %s, want nil", err.Error())
			}

			if !reflect.DeepEqual(tc.wantDevices, devices) {
				t.Errorf("got response body: %+v, want %+v", devices, tc.wantDevices)
			}
		}

		// Test for GET /devices/:id response
		if tc.wantDevice != nil {
			p := tc.params.ByName("id")
			id, _ := strconv.Atoi(p)

			if id != tc.reader.id {
				t.Errorf("got Read() param: %d, want %d", id, tc.reader.id)
			}

			var device *device.Device

			if err := json.NewDecoder(res.Body).Decode(&device); err != nil {
				t.Errorf("got error while decoding response body: %s, want nil", err.Error())
			}

			if !reflect.DeepEqual(tc.wantDevice, device) {
				t.Errorf("got response body: %+v, want %+v", device, tc.wantDevice)
			}
		}
	}
}

func TestGetDevicesError(t *testing.T) {
	testCases := []struct {
		path            string
		params          httprouter.Params
		getter          *getterTest
		reader          *deviceReader
		wantError       Err
		wantContentType string
		wantStatusCode  int
	}{
		{
			"/devices",
			httprouter.Params{},
			&getterTest{},
			&deviceReader{
				err: errors.New("Testing Error"),
			},
			Err{"Testing Error"},
			"application/json",
			http.StatusInternalServerError,
		},
		{
			"/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			&getterTest{},
			&deviceReader{
				devices: []*device.Device{},
			},
			Err{http.StatusText(http.StatusNotFound)},
			"application/json",
			http.StatusNotFound,
		},
		{
			"/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "xzy"},
			},
			&getterTest{},
			&deviceReader{
				devices: []*device.Device{},
			},
			Err{errInvalidID.Error()},
			"application/json",
			http.StatusBadRequest,
		},
		{
			"/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "0"},
			},
			&getterTest{},
			&deviceReader{
				devices: []*device.Device{},
			},
			Err{errInvalidID.Error()},
			"application/json",
			http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		GetDevices(tc.getter, tc.reader)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)

		tc.getter.handle[tc.path](rw, req, tc.params)

		res := rw.Result()

		if res.StatusCode != tc.wantStatusCode {
			t.Errorf("got http status code %d, want %d", res.StatusCode, tc.wantStatusCode)
		}

		if c := rw.Header().Get("Content-Type"); c != tc.wantContentType {
			t.Errorf("got http header `Content-Type`: %s, want %s", c, tc.wantContentType)
		}

		var e Err

		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			t.Errorf("got error while decoding response body: %s, want nil", err.Error())
		}

		if !reflect.DeepEqual(tc.wantError, e) {
			t.Errorf("got response body: %+v, want %+v", e, tc.wantError)
		}

	}
}
