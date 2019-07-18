package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/listing/device"
)

var errInvalidID = errors.New(`id parameter must be an integer and greater than 0`)

// Err is the struct that will be returned as
// the responste to a failed HTTP request
type Err struct {
	Error string
}

// HTTPGetter is an interface that handles the datails of
// parsing an GET request to the server
type HTTPGetter interface {
	GET(path string, h httprouter.Handle)
}

// Devices is the struct that will be returned as
// the response to a succesful HTTP request for listing
// all devices
type Devices struct {
	Data []*device.Device
}

// GetDevices handles the requests for reading devices
func GetDevices(router HTTPGetter, r device.Reader) {
	router.GET("/devices", getDevicesHandler(r))
	router.GET("/devices/:id", getDevicesHandler(r))
}

func getDevicesHandler(r device.Reader) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var (
			id      int
			output  interface{}
			encoder = json.NewEncoder(rw)
		)

		rw.Header().Set("Content-Type", "application/json")

		if param := p.ByName("id"); param != "" {
			pid, err := strconv.Atoi(param)
			if err != nil || pid == 0 {
				rw.WriteHeader(http.StatusBadRequest)
				encoder.Encode(Err{errInvalidID.Error()})
				return
			}

			id = pid
		}

		devices, err := r.Read(id)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(Err{err.Error()})
			return
		}

		if id > 0 && len(devices) == 0 {
			rw.WriteHeader(http.StatusNotFound)
			encoder.Encode(Err{http.StatusText(http.StatusNotFound)})
			return
		}

		rw.WriteHeader(http.StatusOK)

		if id > 0 {
			output = devices[0]
		} else {
			output = Devices{devices}
		}

		encoder.Encode(output)
	}
}
