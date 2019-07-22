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
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/removing/device"
	deviceList "github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/listing/device"
)

type removerTest struct {
	*routerTest
}

func (r *removerTest) DELETE(path string, h httprouter.Handle) {
	r.handleReq(path, h)
}

type deviceRemover struct {
	dvice  *device.Device
	err error
}

func (r *deviceRemover) Remove(d *device.Device) error {
	r.dvice = d
	return r.err
}

func TestRemoveDevice(t *testing.T) {
	testCases := []struct {
		dvices	       []*deviceList.Device
		path           string
		params         httprouter.Params
		wantStatusCode int
	}{
		{
			[]*deviceList.Device{
				&deviceList.Device{ID: 1},
			},
			"/api/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			http.StatusNoContent,
		},
		{
			[]*deviceList.Device{},
			"/api/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		remover := &deviceRemover{}
		reader := &deviceReader{devices: tc.dvices}

		deleter := &removerTest{&routerTest{}}
		RemoveDevice(deleter, remover, reader)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, tc.path, nil)

		deleter.handle[tc.path](rw, req, tc.params)

		res := rw.Result()

		if res.StatusCode != tc.wantStatusCode {
			t.Errorf("got http status code %d, want %d", res.StatusCode, tc.wantStatusCode)
		}

		if id, _ := strconv.Atoi(tc.params.ByName("id")); reader.id != id {
			t.Errorf("got id sent to Read(id): %d, want %d", reader.id, id)
		}

		if len(reader.devices) > 0 {
			if id := reader.devices[0].ID; remover.dvice.ID != id {
				t.Errorf("got id of remover.Device: %d, want %d", remover.dvice.ID, id)
			}

			if user := reader.devices[0].User; remover.dvice.User != user {
				t.Errorf("got id of remover.Device: %d, want %d", remover.dvice.User, user)
			}
		}
	}
}

func TestRemoveDeviceError(t *testing.T) {
	testCases := []struct {
		path            string
		params          httprouter.Params
		reader		*deviceReader
		remover         *deviceRemover
		wantError       Err
		wantStatusCode  int
		wantContentType string
	}{
		{
			"/api/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			&deviceReader{
				err: errors.New("Reading Error"),
			},
			&deviceRemover{},
			Err{"Reading Error"},
			http.StatusInternalServerError,
			"application/json",
		},
		{
			"/api/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			&deviceReader{
				devices: []*deviceList.Device{
					&deviceList.Device{ID: 1},
				},
			},
			&deviceRemover{
				err: &device.InvalidOperation{"Invalid Op"},
			},
			Err{"Invalid Op"},
			http.StatusBadRequest,
			"application/json",
		},
		{
			"/api/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			&deviceReader{
				devices: []*deviceList.Device{
					&deviceList.Device{ID: 1},
				},
			},
			&deviceRemover{
				err: errors.New("Testing Error"),
			},
			Err{"Testing Error"},
			http.StatusInternalServerError,
			"application/json",
		},
		{
			"/api/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "xyz"},
			},
			&deviceReader{
				devices: []*deviceList.Device{
					&deviceList.Device{ID: 1},
				},
			},
			&deviceRemover{},
			Err{errInvalidID.Error()},
			http.StatusBadRequest,
			"application/json",
		},
	}

	for _, tc := range testCases {
		deleter := &removerTest{&routerTest{}}
		RemoveDevice(deleter, tc.remover, tc.reader)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, tc.path, nil)

		deleter.handle[tc.path](rw, req, tc.params)

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
