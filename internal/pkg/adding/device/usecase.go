package device

// Writer is the interface that writes the Device in
// Persistency layer
type Writer interface {
	// Returns the ID the new Device
	Write(d *Device) (ID int, err error)
}
