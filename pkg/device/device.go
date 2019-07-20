package device

import (
	"github.com/jinzhu/gorm"

	deviceWriter "github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/adding/device"
	deviceReader "github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/listing/device"
	//deviceRemover "github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/removing/device"
	deviceUpdater "github.com/rafaeldias/mobilitee-back-end-developer/internal/pkg/updating/device"
)

// Device is the struct that exposes the device operations
// to public use
type Device struct {
	deviceReader.Reader
	deviceWriter.Writer
	deviceUpdater.Updater
	//deviceRemover.Remover
}

// New returns a public Device
func New(db *gorm.DB) *Device {
	return &Device{
		deviceReader.New(deviceReader.NewRepository(db)),
		deviceWriter.New(deviceWriter.NewRepository(db)),
		deviceUpdater.New(deviceUpdater.NewRepository(db)),
		//deviceRemover.New(deviceRemover.NewRepository(db)),
	}
}
