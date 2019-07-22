package api

import (
	"testing"

	"github.com/rafaeldias/mobilitee-back-end-developer/pkg/device"
	"github.com/julienschmidt/httprouter"
)

type routerTest struct {
	getInvoked    bool
	postInvoked   bool
	patchInvoked  bool
	deleteInvoked bool
}

func (r *routerTest) GET(path string, h httprouter.Handle) {
	r.getInvoked = true
}

func (r *routerTest) POST(path string, h httprouter.Handle) {
	r.postInvoked = true
}

func (r *routerTest) PATCH(path string, h httprouter.Handle) {
	r.patchInvoked = true
}

func (r *routerTest) DELETE(path string, h httprouter.Handle) {
	r.deleteInvoked = true
}

func TestRestfulDevice(t *testing.T) {
	rt := &routerTest{}

	RestfulDevice(rt, device.New(nil))

	testCases := []struct {
		method  string
		invoked bool
	}{
		{"GET", rt.getInvoked},
		{"POST", rt.postInvoked},
		{"PATCH", rt.patchInvoked},
		{"DELETE", rt.deleteInvoked},
	}

	for _, tc := range testCases {
		if !tc.invoked {
			t.Errorf("got false for invoking method %s, want true", tc.method)
		}
	}

}
