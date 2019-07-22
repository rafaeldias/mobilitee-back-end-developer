package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/adding/device"
)

// Device is a new device created
type Device struct {
	ID int
}

// HTTPPoster is an interface that handles the datails of
// parsing a POST request to the server
type HTTPPoster interface {
	POST(path string, h httprouter.Handle)
}

// CreateDevice handles the requests for creating devices
func CreateDevice(router HTTPPoster, writer device.Writer) {
	router.POST("/api/devices", createDeviceHandler(writer))
}

func createDeviceHandler(writer device.Writer) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		var (
			dvice   *device.Device
			encoder = json.NewEncoder(rw)
		)

		rw.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(req.Body).Decode(&dvice); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			encoder.Encode(Err{fmt.Sprintf("invalid JSON syntax: %s",
				err.Error())})
			return
		}

		id, err := writer.Write(dvice)
		if err != nil {
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

		rw.Header().Set("Location", fmt.Sprintf("/devices/%d", id))
		rw.WriteHeader(http.StatusCreated)

		encoder.Encode(Device{id})
	}
}
