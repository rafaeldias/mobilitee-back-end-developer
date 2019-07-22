package rest

import (
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
		id, err := strconv.Atoi(param.ByName("id"))
		if err != nil {
			(&Err{errInvalidID.Error()}).Write(rw, http.StatusBadRequest)
			return
		}

		dvices, err := reader.Read(id)
		if err != nil {
			(&Err{err.Error()}).Write(rw, http.StatusInternalServerError)
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

			(&Err{err.Error()}).Write(rw, status)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
