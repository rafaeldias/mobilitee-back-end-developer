package device

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type repoWriteExchanger struct {
	id	    int
	nw	    *Device
	old	    *Device
	exchangeErr error
	writeErr    error
}

func (r *repoWriteExchanger) Write(d *Device) (id int, err error) {
	r.nw = d

	id = r.id
	err = r.writeErr

	return
}

func (r *repoWriteExchanger) Exchange(old, nw *Device) (id int, err error) {
	r.old = old
	r.nw = nw

	id = r.id
	err = r.exchangeErr

	return
}

type testUser struct {
	devices			 int
	device			 int
	latestExchangeExpirestAt time.Time
	exchanging		 bool
	user			 int
	exchangingErr		 error
	latestExchangeErr	 error
	countDevicesErr		 error
}

func (u *testUser) IsExchanging(user int) (exchanging bool, id int, err error) {
	u.user = user

	exchanging = u.exchanging
	id = u.device
	err = u.exchangingErr

	return
}

func (u *testUser) LatestExchangeExpiresAt(user int) (time.Time, error) {
	u.user = user
	return u.latestExchangeExpirestAt, u.latestExchangeErr
}

func (u *testUser) CountDevices(user int) (int, error) {
	u.user = user
	return u.devices, u.countDevicesErr
}

func TestNew(t *testing.T) {
	testCases := []struct {
		repo WriteExchanger
		user User
	}{
		{nil, nil},
		{&repoWriteExchanger{}, &testUser{}},
	}

	for _, tc := range testCases {
		d := New(tc.repo, tc.user)

		if !reflect.DeepEqual(d.repo, tc.repo) {
			t.Errorf("got WriteExchanger: %+v, want %+v", tc.repo, tc.repo)
		}

		if !reflect.DeepEqual(d.user, tc.user) {
			t.Errorf("got User: %+v, want %+v", tc.user, tc.user)
		}
	}
}

func TestUsecaseWrite(t *testing.T) {
	testCases := []struct {
		id    int
		repo  *repoWriteExchanger
		user  *testUser
		dvice *Device
	}{
		{
			1,
			&repoWriteExchanger{
				id: 1,
			},
			&testUser{},
			&Device{Name: "Test", Model: "Android", User: 1},
		},
		{
			1,
			&repoWriteExchanger{
				id: 1,
			},
			&testUser{exchanging: true, device: 3},
			&Device{Name: "Test", Model: "Android", User: 1},
		},
	}

	for _, tc := range testCases {
		device := New(tc.repo, tc.user)

		id, err := device.Write(tc.dvice)
		if err != nil {
			t.Errorf("got error while calling Device.Write(%+v): %s, want nil",
				tc.dvice, err.Error())
		}

		if tc.id != id {
			t.Errorf("got devices from Device.Write(%+v): %d, want: %d", tc.dvice,
				id, tc.id)
		}

		if !reflect.DeepEqual(tc.dvice, tc.repo.nw) {
			t.Errorf("got repo.nw from calling Device.Write(%+v): %+v, want: %+v", tc.dvice,
				tc.dvice, tc.repo.nw)
		}

		if tc.dvice.User != tc.user.user {
			t.Errorf("got user.user from calling Device.Write(%+v): %d, want: %d", tc.dvice,
				tc.user.user, tc.dvice.User)
		}

		if tc.user.exchanging && !reflect.DeepEqual(tc.user.device, tc.repo.old.ID) {
			t.Errorf("got repo.old from calling Device.Write(%+v): %d, want: %d", tc.dvice,
				tc.repo.old.ID, tc.user.device)
		}
	}
}

func TestUsecaseWriteError(t *testing.T) {
	testCases := []struct {
		id    int
		repo  *repoWriteExchanger
		user  *testUser
		dvice *Device
	}{
		{
			1,
			&repoWriteExchanger{},
			&testUser{},
			&Device{Name: "", User: 1},
		},
		{
			1,
			&repoWriteExchanger{},
			&testUser{},
			&Device{Name: "Testing", User: 0},
		},
		{
			1,
			&repoWriteExchanger{},
			&testUser{},
			&Device{Name: "Testing", Model: "", User: 1},
		},
		{
			1,
			&repoWriteExchanger{},
			&testUser{exchangingErr: errors.New("IsExchanging error")},
			&Device{Name: "Testing", Model: "Android", User: 1},
		},
		{
			1,
			&repoWriteExchanger{},
			&testUser{latestExchangeErr: errors.New("LatestExchangeExpiresAt error")},
			&Device{Name: "Testing", Model: "Android", User: 1},
		},
		{
			1,
			&repoWriteExchanger{},
			&testUser{countDevicesErr: errors.New("CountDevices error")},
			&Device{Name: "Testing", Model: "Android", User: 1},
		},
		{
			1,
			&repoWriteExchanger{exchangeErr: errors.New("Exchange Error")},
			&testUser{exchanging: true},
			&Device{Name: "Testing", Model: "Android", User: 1},
		},
		{
			1,
			&repoWriteExchanger{writeErr: errors.New("Write Error")},
			&testUser{},
			&Device{Name: "Testing", Model: "Android", User: 1},
		},
	}

	for _, tc := range testCases {
		device := New(tc.repo, tc.user)

		_, err := device.Write(tc.dvice)
		if err == nil {
			t.Errorf("got error nil while calling Device.Write(%+v), want nil",
				tc.dvice)
		}
	}
}
