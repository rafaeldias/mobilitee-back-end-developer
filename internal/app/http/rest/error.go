package rest

import (
	"encoding/json"
	"errors"
	"net/http"
)

var errInvalidID = errors.New(`id parameter must be an integer and greater than 0`)

// Err is the struct that will be returned as
// the responste to a failed HTTP request
type Err struct {
	Error string
}

// Write writes the error with it code to the ResponseWriter
func (e *Err) Write(rw http.ResponseWriter, code int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(e)
}

