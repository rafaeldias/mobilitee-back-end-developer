package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/updating/device"
)

// HTTPPatch is an interface that handles the datails of
// parsing an PATCH request to the server
type HTTPPatch interface {
	PATCH(path string, h httprouter.Handle)
}

// UpdateDevice handles the requests for updating devices
func UpdateDevice(router HTTPPatch, updater device.Updater) {
	router.PATCH("/devices/:id", updateDeviceHandler(updater))
}

func updateDeviceHandler(updater device.Updater) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, param httprouter.Params) {
		var (
			dvice   *device.Device
			encoder = json.NewEncoder(rw)
		)

		if err := json.NewDecoder(req.Body).Decode(&dvice); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Err{fmt.Sprintf("invalid JSON syntax: %s",
				err.Error())})
			return
		}

		id, err := strconv.Atoi(param.ByName("id"))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Err{errInvalidID.Error()})
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

			rw.WriteHeader(status)
			encoder.Encode(Err{err.Error()})
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
