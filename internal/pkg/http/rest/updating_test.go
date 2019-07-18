package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/updating/device"
)

type patcherTest struct {
	*routerTest
}

func (p *patcherTest) PATCH(path string, h httprouter.Handle) {
	p.handleReq(path, h)
}

type deviceUpdater struct {
	id    int
	err   error
	dvice *device.Device
}

func (u *deviceUpdater) Update(id int, d *device.Device) error {
	u.id = id
	u.dvice = d
	return u.err
}

func TestUpdateDevices(t *testing.T) {
	testCases := []struct {
		path           string
		body           string
		params         httprouter.Params
		dvice          *device.Device
		wantStatusCode int
	}{
		{
			"/devices/:id",
			`{"Name":"Testing"}`,
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			&device.Device{Name: "Testing"},
			http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		updater := &deviceUpdater{}

		patcher := &patcherTest{&routerTest{}}
		UpdateDevice(patcher, updater)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, tc.path, bytes.NewBuffer([]byte(tc.body)))

		patcher.handle[tc.path](rw, req, tc.params)

		res := rw.Result()

		if res.StatusCode != tc.wantStatusCode {
			t.Errorf("got http status code %d, want %d", res.StatusCode, tc.wantStatusCode)
		}

		if id, _ := strconv.Atoi(tc.params.ByName("id")); updater.id != id {
			t.Errorf("got id sent to Update(id, device): %d, want %d", updater.id, id)
		}

		if !reflect.DeepEqual(updater.dvice, tc.dvice) {
			t.Errorf("got device sent to Update(id, device): %+v, want %+v",
				updater.dvice, tc.dvice)
		}

	}
}

func TestUpdateDeviceError(t *testing.T) {
	testCases := []struct {
		path           string
		body           string
		params         httprouter.Params
		updater        *deviceUpdater
		wantError      Err
		wantStatusCode int
	}{
		{
			"/devices/:id",
			``,
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			&deviceUpdater{},
			Err{"invalid JSON syntax: EOF"},
			http.StatusBadRequest,
		},
		{
			"/devices/:id",
			`{}`,
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			&deviceUpdater{
				err: errors.New("Testing Error"),
			},
			Err{"Testing Error"},
			http.StatusInternalServerError,
		},
		{
			"/devices/:id",
			`{"Name": ""}`,
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			&deviceUpdater{
				err: &device.ErrInvalidInput{"ID", "empty"},
			},
			Err{(&device.ErrInvalidInput{"ID", "empty"}).Error()},
			http.StatusBadRequest,
		},
		{
			"/devices/:id",
			`{"Name": ""}`,
			httprouter.Params{
				httprouter.Param{"id", "2"},
			},
			&deviceUpdater{
				err: &device.ErrNotFound{2},
			},
			Err{(&device.ErrNotFound{2}).Error()},
			http.StatusNotFound,
		},
		{
			"/devices/:id",
			`{"Name": ""}`,
			httprouter.Params{
				httprouter.Param{"id", "xyz"},
			},
			&deviceUpdater{},
			Err{errInvalidID.Error()},
			http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		patcher := &patcherTest{&routerTest{}}
		UpdateDevice(patcher, tc.updater)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, tc.path, bytes.NewBuffer([]byte(tc.body)))

		patcher.handle[tc.path](rw, req, tc.params)

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
	}
}
