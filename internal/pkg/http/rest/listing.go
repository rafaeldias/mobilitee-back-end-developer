package rest

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/listing/device"
)

// HTTPGetter is an interface that handles the datails of
// parsing an GET request to the server
type HTTPGetter interface {
	GET(path string, h httprouter.Handle)
}

// GetDevices handles the requests for reading devices
func GetDevices(router HTTPGetter, r device.Reader) {
	router.GET("/devices", func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		rw.WriteHeader(http.StatusOK)
	})
	router.GET("/devices/:id", nil)
}
