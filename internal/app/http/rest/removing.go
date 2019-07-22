package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/removing/device"
	deviceList "github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/listing/device"
)

// HTTPDeleter is an interface that handles the datails of
// parsing an DELETE request to the server
type HTTPDeleter interface {
	DELETE(path string, h httprouter.Handle)
}

// RemoveDevice handles the requests for removing devices
func RemoveDevice(router HTTPDeleter, remover device.Remover, reader deviceList.Reader) {
	router.DELETE("/api/devices/:id", removeDeviceHandler(remover, reader))
}

func removeDeviceHandler(remover device.Remover, reader deviceList.Reader) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, param httprouter.Params) {
		var (
			encoder = json.NewEncoder(rw)
			header  = rw.Header()
		)

		id, err := strconv.Atoi(param.ByName("id"))
		if err != nil {
			header.Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Err{errInvalidID.Error()})
			return
		}

		dvices, err := reader.Read(id)
		if err != nil {
			header.Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(Err{err.Error()})
			return
		}

		if len(dvices) == 0 {
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		dvice := &device.Device{ID: dvices[0].ID, User: dvices[0].User}

		if err := remover.Remove(dvice); err != nil {
			var status int

			switch err.(type) {
			case *device.InvalidOperation:
				status = http.StatusBadRequest
			default:
				status = http.StatusInternalServerError
			}

			header.Set("Content-Type", "application/json")
			rw.WriteHeader(status)
			encoder.Encode(Err{err.Error()})
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
