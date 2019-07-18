package dvice

// Remover is the interface that removes
// the device from pesistency layer
type Remover interface {
	Remove(id int) error
}
