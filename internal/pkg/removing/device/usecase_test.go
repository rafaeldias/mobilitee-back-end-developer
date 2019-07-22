package device

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type repoRemover struct {
	dvice  *Device
	err	error
}

func (r *repoRemover) Remove(d *Device) error {
	r.dvice = d
	return r.err
}

type testUser struct {
	devices			 int
	latestExchangeExpirestAt time.Time
	id			 int
	latestExchangeErr	 error
	countDevicesErr		 error
}

func (u *testUser) LatestExchangeExpiresAt(user int) (time.Time, error) {
	u.id = user
	return u.latestExchangeExpirestAt, u.latestExchangeErr
}

func (u *testUser) CountDevices(user int) (int, error) {
	u.id = user
	return u.devices, u.countDevicesErr
}

func TestNew(t *testing.T) {
	testCases := []struct {
		repo Remover
		user User
	}{
		{nil, nil},
		{&repoRemover{}, &testUser{}},
	}

	for _, tc := range testCases {
		d := New(tc.repo, tc.user)

		if !reflect.DeepEqual(d.repo, tc.repo) {
			t.Errorf("got Remover: %+v, want %+v", tc.repo, tc.repo)
		}

		if !reflect.DeepEqual(d.user, tc.user) {
			t.Errorf("got User: %+v, want %+v", tc.user, tc.user)
		}
	}
}

func TestUsecaseRemove(t *testing.T) {
	testCases := []struct {
		id    int
		repo  *repoRemover
		dvice *Device
	}{
		{
			1,
			&repoRemover{},
			&Device{ID: 1, User: 1},
		},
	}

	for _, tc := range testCases {
		device := New(tc.repo, &testUser{})

		err := device.Remove(tc.dvice)
		if err != nil {
			t.Fatalf("got error while calling Usecase.Remove(%+v): %s, want nil",
				tc.dvice, err.Error())
		}

		if !reflect.DeepEqual(tc.dvice, tc.repo.dvice) {
			t.Errorf("got repo.dvice from calling Usecae.Remove(%+v): %+v, want: %+v", tc.dvice,
				tc.repo.dvice, tc.dvice)
		}
	}
}

func TestUsecaseRemoveError(t *testing.T) {
	testCases := []struct {
		repo  *repoRemover
		user  *testUser
		dvice *Device
	}{
		{
			&repoRemover{},
			&testUser{
				devices : 1,
				latestExchangeExpirestAt: time.Now().AddDate(0, 0, -10),
			},
			&Device{ID: 2, User: 1},
		},
		{
			&repoRemover{},
			&testUser{
				latestExchangeErr: errors.New("Lastest Exchange Error"),
			},
			&Device{ID: 2, User: 1},
		},
		{
			&repoRemover{},
			&testUser{
				countDevicesErr: errors.New("Counting Devices Error"),
			},
			&Device{ID: 2, User: 1},
		},
	}

	for _, tc := range testCases {
		device := New(tc.repo, tc.user)

		if err := device.Remove(tc.dvice); err == nil {
			t.Errorf("got error nil while calling Usecase.Remove(%+v), want not nil",
				tc.dvice)
		}

		if tc.user.id != tc.dvice.User {
			t.Errorf("got user.user while calling Usecase.Remove(%+v): %d; want %d",
				tc.dvice, tc.user.id, tc.dvice.User)
		}
	}
}
