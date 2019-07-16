package device

// Reader is an interface for reading a list of Devices.
// It returns a list of devices or error
type Reader interface {
	Read(ID int) ([]*Device, error)
}

// Usecase is a application and business specific layer
// for orchastrating the flow of data to and from external
// channels (http, database, etc.)
type Usecase struct {
	// repo is the presistency layer, which must implement
	// the Reader interface
	repo Reader
}

// New returns a new Usecase
func New(repo Reader) *Usecase {
	return &Usecase{repo}
}

// Read implements the Reader interface for listing devices
func (u *Usecase) Read(ID int) ([]*Device, error) {
	return u.repo.Read(ID)
}
