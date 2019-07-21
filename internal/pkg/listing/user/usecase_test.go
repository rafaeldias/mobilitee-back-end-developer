package user

import (
	"errors"
	"testing"
	"reflect"
	"time"
)

type repoDeviceReader struct {
	dvice *Device
	err error
	user int
	dvices int
}

func (r *repoDeviceReader) LatestRemoved(user int) (*Device, error) {
	r.user = user
	return r.dvice, r.err
}

func (r *repoDeviceReader) LatestExchange(user int) (*Device, error) {
	r.user = user
	return r.dvice, r.err
}

func (r *repoDeviceReader) Count(user int) (int, error) {
	r.user = user
	return r.dvices, r.err
}

func TestNew(t *testing.T) {
	testCases := []struct {
		repo *repoDeviceReader
	}{
		{nil},
		{&repoDeviceReader{}},
	}

	for _, tc := range testCases {
		d := New(tc.repo)

		if !reflect.DeepEqual(d.repo, tc.repo) {
			t.Errorf("got New(%+v): %+v, want: %+v", tc.repo, d.repo, tc.repo)
		}
	}
}

func TestUsecaseIsExchanging(t *testing.T) {
	testCases := []struct {
		dvice		*Device
		wantExchanging  bool
		wantDeviceID    int
		wantUser	int
	}{
		{
			&Device{ID: 1, LatestRemovedAt: time.Now().AddDate(0, 0, -20)},
			true,
			1,
			20,
		},
		{
			&Device{ID: 1, LatestRemovedAt: time.Now().AddDate(0, 0, -31)},
			false,
			0,
			20,
		},
	}

	for _, tc := range testCases {
		repo := &repoDeviceReader{dvice: tc.dvice}

		u := New(repo)

		exchanging, dvice, err := u.IsExchanging(tc.wantUser)
		if err != nil {
			t.Fatalf("got error while calling Usecase.IsExchanging(%d): %s, want nil",
				tc.wantUser, err.Error())
		}

		if exchanging != tc.wantExchanging {
			t.Errorf("got exchanging %t, want %t", exchanging, tc.wantExchanging)
		}

		if dvice != tc.wantDeviceID {
			t.Errorf("got device ID %d, want %d", dvice, tc.wantDeviceID)
		}

		if repo.user != tc.wantUser {
			t.Errorf("got repo.user %d, want %d", repo.user, tc.wantUser)
		}

	}
}

func TestUsecaseIsExchangingError(t *testing.T) {
		repo := &repoDeviceReader{err: errors.New("Testing Error")}

		u := New(repo)

		_, _, err := u.IsExchanging(1)
		if err == nil {
			t.Error("got error nil while calling Usecase.IsExchanging(1), want not nil")
		}
}

func TestUsecaseLatestExchangeExpiresAt(t *testing.T) {
	testCases := []struct {
		dvice		     *Device
		wantNextExchangingAt time.Time
		wantUser	     int
	}{
		{
			&Device{LatestExchangeAt: time.Now().AddDate(0, 0, -10)},
			time.Now().Add(time.Hour * 24 * 30 - time.Since(time.Now().AddDate(0, 0, -10))),
			1,
		},
		{
			&Device{LatestExchangeAt: time.Now().AddDate(0, 0, -31)},
			time.Time{},
			1,
		},
	}

	for _, tc := range testCases {
		repo := &repoDeviceReader{dvice: tc.dvice}

		u := New(repo)

		nextExchangingAt, err := u.LatestExchangeExpiresAt(tc.wantUser)
		if err != nil {
			t.Fatalf("got error while calling Usecase.LatestExchangeExpiresAt(%d): %s, want nil",
				tc.wantUser, err.Error())
		}

		if !reflect.DeepEqual(tc.wantNextExchangingAt, nextExchangingAt) {
			t.Errorf("got nextExchangingAt %s, want %s", nextExchangingAt, tc.wantNextExchangingAt)
		}

		if repo.user != tc.wantUser {
			t.Errorf("got repo.user %d, want %d", repo.user, tc.wantUser)
		}

	}
}

func TestUsecaseLatestExchangeExpiresAtError(t *testing.T) {
		repo := &repoDeviceReader{err: errors.New("Testing Error")}

		u := New(repo)

		_, err := u.LatestExchangeExpiresAt(1)
		if err == nil {
			t.Error("got error nil while calling Usecase.LatestExchangeExpiresAt(1), want not nil")
		}
}

func TestUsecaseCountDevices(t *testing.T) {
		repo := &repoDeviceReader{dvices: 2}
		want := 2
		wantUser := 1

		u := New(repo)

		got, err := u.CountDevices(wantUser)
		if err != nil {
		t.Fatalf("got error while calling Usecase.CountDevices(%d): %s, want nil",
				wantUser, err.Error())
		}

		if got != want {
			t.Errorf("got number of devices %d, want %d", got, want)
		}

		if repo.user != wantUser {
			t.Errorf("got repo.user %d, want %d", repo.user, wantUser)
		}
}

func TestUsecaseCountDevicesError(t *testing.T) {
		repo := &repoDeviceReader{err: errors.New("Testing Error")}

		u := New(repo)

		_, err := u.CountDevices(1)
		if err == nil {
			t.Fatalf("got error nil while calling Usecase.CountDevices(1), want not nil")
		}
}
