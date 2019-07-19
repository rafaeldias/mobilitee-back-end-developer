package api

import (
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/http/rest"
	"github.com/rafaeldias/mobilitee-back-end-developer/pkg/device"
)

// Router handle the HTTP methods for the rest API
type Router interface {
	rest.HTTPGetter
	rest.HTTPPoster
	rest.HTTPPatcher
	rest.HTTPDeleter
}

// RestfulDevice setup the REST api handlers for device
func RestfulDevice(r Router, d *device.Device) {
	rest.GetDevices(r, d)
	//rest.CreateDevice(r, d)
	rest.UpdateDevice(r, d)
	//rest.RemoveDevice(r, d)
}
