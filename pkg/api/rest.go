package api

import (
	"github.com/rafaeldias/mobilitee-back-end-developer/internal/app/http/rest"
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
	rest.GetDevices(r, d.Reader)
	rest.CreateDevice(r, d.Writer)
	rest.UpdateDevice(r, d.Updater)
	rest.RemoveDevice(r, d.Remover, d.Reader)
}
