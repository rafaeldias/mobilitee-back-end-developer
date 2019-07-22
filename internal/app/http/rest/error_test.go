package rest

import (
	"encoding/json"
	"testing"
	"net/http/httptest"
	"net/http"
)

func TestErrorWrite(t *testing.T) {
	rw := httptest.NewRecorder()

	err := &Err{"Testing Error"}

	err.Write(rw, http.StatusBadRequest)

	if c := rw.Header().Get("Content-Type"); c != "application/json" {
		t.Errorf("got Header `Content-Type`: %s, want application/json", c)
	}

	if rw.Code != http.StatusBadRequest {
		t.Errorf("got Status Code: %d, want %d", rw.Code, http.StatusBadRequest)
	}

	var e Err
	if err := json.NewDecoder(rw.Body).Decode(&e); err != nil {
		t.Fatalf("got error while decoding response body: %s, want nil", err.Error())
	}

	if e.Error != err.Error {
		t.Errorf("got Error Message: %s, want %s", e.Error, err.Error)
	}
}
