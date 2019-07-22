package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/updating/device"
)

// HTTPPatcher is an interface that handles the datails of
// parsing an PATCH request to the server
type HTTPPatcher interface {
	PATCH(path string, h httprouter.Handle)
}

// UpdateDevice handles the requests for updating devices
func UpdateDevice(router HTTPPatcher, updater device.Updater) {
	router.PATCH("/api/devices/:id", updateDeviceHandler(updater))
}

func updateDeviceHandler(updater device.Updater) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, param httprouter.Params) {
		var (
			dvice   *device.Device
		)

		if err := json.NewDecoder(req.Body).Decode(&dvice); err != nil {
			(&Err{fmt.Sprintf("invalid JSON syntax: %s",
				err.Error())}).Write(rw, http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(param.ByName("id"))
		if err != nil {
			(&Err{errInvalidID.Error()}).Write(rw, http.StatusBadRequest)
			return
		}

		if err := updater.Update(id, dvice); err != nil {
			var status int

			switch err.(type) {
			case *device.ErrInvalidInput:
				status = http.StatusBadRequest
			case *device.ErrNotFound:
				status = http.StatusNotFound
			default:
				status = http.StatusInternalServerError
			}

			(&Err{err.Error()}).Write(rw, status)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
