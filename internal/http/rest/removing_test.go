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
	//"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/updating/device"
)

type removerTest struct {
	*routerTest
}

func (r *removerTest) DELETE(path string, h httprouter.Handle) {
	r.handleReq(path, h)
}

type deviceRemover struct {
	id  int
	err error
}

func (r *deviceRemover) Remove(id int) error {
	r.id = id
	return r.err
}

func TestRemoveDevices(t *testing.T) {
	testCases := []struct {
		path           string
		params         httprouter.Params
		wantStatusCode int
	}{
		{
			"/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		remover := &deviceRemover{}

		deleter := &removerTest{&routerTest{}}
		RemoveDevice(deleter, remover)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, tc.path, nil)

		deleter.handle[tc.path](rw, req, tc.params)

		res := rw.Result()

		if res.StatusCode != tc.wantStatusCode {
			t.Errorf("got http status code %d, want %d", res.StatusCode, tc.wantStatusCode)
		}

		if id, _ := strconv.Atoi(tc.params.ByName("id")); remover.id != id {
			t.Errorf("got id sent to Remove(id, device): %d, want %d", remover.id, id)
		}
	}
}

func TestRemoveDeviceError(t *testing.T) {
	testCases := []struct {
		path           string
		params         httprouter.Params
		remover        *deviceRemover
		wantError      Err
		wantStatusCode int
	}{
		{
			"/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "1"},
			},
			&deviceRemover{
				err: errors.New("Testing Error"),
			},
			Err{"Testing Error"},
			http.StatusInternalServerError,
		},
		{
			"/devices/:id",
			httprouter.Params{
				httprouter.Param{"id", "xyz"},
			},
			&deviceRemover{},
			Err{errInvalidID.Error()},
			http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		deleter := &removerTest{&routerTest{}}
		RemoveDevice(deleter, tc.remover)

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
	}
}
