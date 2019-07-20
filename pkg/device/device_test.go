package device

import "testing"

var errMsg = "got Device.%s nil, want not nil"

func TestNew(t *testing.T) {
	d := New(nil)

	if d.Reader == nil {
		t.Errorf(errMsg, "Reader")
	}

	if d.Writer == nil {
		t.Errorf(errMsg, "Writer")
	}

	if d.Updater == nil {
		t.Errorf(errMsg, "Updater")
	}

	//if d.Remover == nil {
	//	t.Errorf(errMsg, "Remover")
	//}

}
