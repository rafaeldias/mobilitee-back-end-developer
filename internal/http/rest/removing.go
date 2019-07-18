package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/removing/device"
)

// HTTPRemover is an interface that handles the datails of
// parsing an DELETE request to the server
type HTTPRemover interface {
	DELETE(path string, h httprouter.Handle)
}

// RemoveDevice handles the requests for removing devices
func RemoveDevice(router HTTPRemover, remover device.Remover) {
	router.DELETE("/devices/:id", removeDeviceHandler(remover))
}

func removeDeviceHandler(remover device.Remover) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, param httprouter.Params) {
		var encoder = json.NewEncoder(rw)

		id, err := strconv.Atoi(param.ByName("id"))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Err{errInvalidID.Error()})
			return
		}

		if err := remover.Remove(id); err != nil {
			var status int

			switch err.(type) {
			// TODO: include validation errors here
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
