package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/adding/device"
)

type posterTest struct {
	*routerTest
}

func (p *posterTest) POST(path string, h httprouter.Handle) {
	p.handleReq(path, h)
}

type deviceWriter struct {
	id    int
	err   error
	dvice *device.Device
}

func (w *deviceWriter) Write(d *device.Device) (id int, err error) {
	w.dvice = d

	id = w.id
	err = w.err
	return
}

func TestCreateDevice(t *testing.T) {
	testCases := []struct {
		path               string
		body               string
		writer             *deviceWriter
		newDevice          *device.Device
		wantDevice         *Device
		wantStatusCode     int
		wantLocationHeader string
		wantContentType    string
	}{
		{
			"/api/devices",
			`{"Name":"Testing","Model":"Android","User":1}`,
			&deviceWriter{
				id: 1,
			},
			&device.Device{
				Name:  "Testing",
				Model: "Android",
				User:  1,
			},
			&Device{ID: 1},
			http.StatusCreated,
			"/devices/1",
			"application/json",
		},
	}

	for _, tc := range testCases {
		poster := &posterTest{&routerTest{}}
		CreateDevice(poster, tc.writer)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, tc.path, bytes.NewBuffer([]byte(tc.body)))

		poster.handle[tc.path](rw, req, httprouter.Params{})

		res := rw.Result()

		if res.StatusCode != tc.wantStatusCode {
			t.Errorf("got http status code %d, want %d", res.StatusCode, tc.wantStatusCode)
		}

		if !reflect.DeepEqual(tc.writer.dvice, tc.newDevice) {
			t.Errorf("got device sent to Write(device): %+v, want %+v",
				tc.writer.dvice, tc.newDevice)
		}

		var dvice *Device

		if err := json.NewDecoder(res.Body).Decode(&dvice); err != nil {
			t.Errorf("got error while decoding response body: %s, want nil", err.Error())
		}

		if !reflect.DeepEqual(tc.wantDevice, dvice) {
			t.Errorf("got response body: %+v, want %+v", dvice, tc.wantDevice)
		}

		if c := rw.Header().Get("Location"); c != tc.wantLocationHeader {
			t.Errorf("got http header `Location`: %s, want %s", c, tc.wantLocationHeader)
		}

		if c := rw.Header().Get("Content-Type"); c != tc.wantContentType {
			t.Errorf("got http header `Content-Type`: %s, want %s", c, tc.wantContentType)
		}
	}
}

func TestCreateDeviceError(t *testing.T) {
	testCases := []struct {
		path            string
		body            string
		writer          *deviceWriter
		wantError       Err
		wantStatusCode  int
		wantContentType string
	}{
		{
			"/api/devices",
			``,
			&deviceWriter{},
			Err{"invalid JSON syntax: EOF"},
			http.StatusBadRequest,
			"application/json",
		},
		{
			"/api/devices",
			`{}`,
			&deviceWriter{
				err: errors.New("Testing Error"),
			},
			Err{"Testing Error"},
			http.StatusInternalServerError,
			"application/json",
		},
	}

	for _, tc := range testCases {
		poster := &posterTest{&routerTest{}}
		CreateDevice(poster, tc.writer)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, tc.path, bytes.NewBuffer([]byte(tc.body)))

		poster.handle[tc.path](rw, req, httprouter.Params{})

		res := rw.Result()

		if res.StatusCode != tc.wantStatusCode {
			t.Errorf("got http status code %d, want %d", res.StatusCode, tc.wantStatusCode)
		}

		var e Err

		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			t.Errorf("got error while decoding response body: %s, want nil", err.Error())
		}

		if !reflect.DeepEqual(tc.wantError, e) {
			t.Errorf("got response body: %+v, want %+v", e, tc.wantError)
		}

		if c := rw.Header().Get("Content-Type"); c != tc.wantContentType {
			t.Errorf("got http header `Content-Type`: %s, want %s", c, tc.wantContentType)
		}
	}
}
