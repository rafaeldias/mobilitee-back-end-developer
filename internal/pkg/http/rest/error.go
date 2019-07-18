package rest

import "errors"

var errInvalidID = errors.New(`id parameter must be an integer and greater than 0`)

// Err is the struct that will be returned as
// the responste to a failed HTTP request
type Err struct {
	Error string
}
