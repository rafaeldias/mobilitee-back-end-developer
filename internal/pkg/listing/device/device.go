package device

import "time"

// Device is the user's device for watching movies and series.
// Used for listing
type Device struct {
	ID        int
	Name      string
	Model     string
	User      int
	CreatedAt time.Time
	DeletedAt time.Time `json:"-"`
}
