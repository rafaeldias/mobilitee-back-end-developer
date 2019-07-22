package user

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func dbSetup(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error while calling sqlmock.New(): %s, want nil", err.Error())
	}

	g, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("got error while calling gorm.Open(\"mysql\", %+v): %s, want nil",
			db, err.Error())
	}

	return g, mock
}

func TestNewRepository(t *testing.T) {
	db, _ := dbSetup(t)

	defer db.Close()

	repo := NewRepository(db)

	if !reflect.DeepEqual(repo.db, db) {
		t.Errorf("got db property of Repository: %+v, want: %+v", repo.db, db)
	}
}

func TestLatestRemoved(t *testing.T) {
	testCases := []struct {
		user          int
		latestRemoved time.Time
		wantDevice    *Device
		wantQuery     string
	}{
		{
			1,
			time.Now(),
			&Device{ID: 1, LatestRemovedAt: time.Now()},
			"SELECT id, MAX\\(deleted_at\\) as latest_removed_at FROM `devices` WHERE \\(\"user\" = \\? AND exchanged IS NULL AND deleted_at IS NOT NULL\\) GROUP BY id ORDER BY deleted_at ASC LIMIT 1",
		},
	}

	for _, tc := range testCases {
		db, mock := dbSetup(t)

		rows := sqlmock.NewRows([]string{"id", "latest_removed_at"}).
				AddRow(tc.wantDevice.ID, tc.latestRemoved)
		mock.ExpectQuery(tc.wantQuery).
			WithArgs(tc.user).
			WillReturnRows(rows)

		repo := NewRepository(db)

		dvice, err := repo.LatestRemoved(tc.user)
		if err != nil {
			t.Errorf("got error while calling Repository.LatestRemoved(%d): %s, want nil",
				tc.user, err.Error())
		}

		if !reflect.DeepEqual(dvice, tc.wantDevice) {
			t.Errorf("go *Device while calling Repository.LatestRemoved(%d): %+v, want %+v",
				tc.user, dvice, tc.wantDevice)
		}

		db.Close()
	}
}

func TestLatestRemovedError(t *testing.T) {
	db, mock := dbSetup(t)
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("Database Error"))

	repo := NewRepository(db)

	_, err := repo.LatestRemoved(1)
	if err == nil {
		t.Error("got error nil while calling Repository.LatestRemoved(1): nil; want not nil")
	}

}

func TestLatestExchange(t *testing.T) {
	testCases := []struct {
		user             int
		latestExchangeAt time.Time
		wantDevice       *Device
		wantQuery        string
	}{
		{
			1,
			time.Now(),
			&Device{LatestExchangeAt: time.Now()},
			"SELECT MAX\\(devices\\.created_at\\) as latest_exchange_at FROM `devices` JOIN devices d ON d\\.exchanged = devices\\.id WHERE \\(devices.user = \\?\\)",
		},
	}

	for _, tc := range testCases {
		db, mock := dbSetup(t)

		rows := sqlmock.NewRows([]string{"latest_exchange_at"}).
				AddRow(tc.latestExchangeAt)
		mock.ExpectQuery(tc.wantQuery).
			WithArgs(tc.user).
			WillReturnRows(rows)

		repo := NewRepository(db)

		dvice, err := repo.LatestExchange(tc.user)
		if err != nil {
			t.Errorf("got error while calling Repository.LatestExchange(%d): %s, want nil",
				tc.user, err.Error())
		}

		if !reflect.DeepEqual(dvice, tc.wantDevice) {
			t.Errorf("go *Device while calling Repository.LatestExchange(%d): %+v, want %+v",
				tc.user, dvice, tc.wantDevice)
		}

		db.Close()
	}
}

func TestLatestExchangeError(t *testing.T) {
	db, mock := dbSetup(t)
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("Database Error"))

	repo := NewRepository(db)

	_, err := repo.LatestExchange(1)
	if err == nil {
		t.Error("got error nil while calling Repository.LatestRemoved(1): nil; want not nil")
	}

}

func TestCount(t *testing.T) {
	testCases := []struct {
		user             int
		count		 int
		wantQuery        string
	}{
		{
			1,
			3,
			"SELECT count\\(\\*\\) FROM `devices` WHERE \\(\"user\" = \\?\\ AND deleted_at IS NULL\\)",
		},
	}

	for _, tc := range testCases {
		db, mock := dbSetup(t)

		rows := sqlmock.NewRows([]string{"count"}).
				AddRow(tc.count)
		mock.ExpectQuery(tc.wantQuery).
			WithArgs(tc.user).
			WillReturnRows(rows)

		repo := NewRepository(db)

		count, err := repo.Count(tc.user)
		if err != nil {
			t.Errorf("got error while calling Repository.LatestExchange(%d): %s, want nil",
				tc.user, err.Error())
		}

		if count != tc.count {
			t.Errorf("go count while calling Repository.Count(%d): %d, want %d",
				tc.user, count, tc.count)
		}

		db.Close()
	}
}

func TestCountError(t *testing.T) {
	db, mock := dbSetup(t)
	defer db.Close()

	mock.ExpectQuery("count").WillReturnError(errors.New("Database Error"))

	repo := NewRepository(db)

	_, err := repo.Count(1)
	if err == nil {
		t.Error("got error nil while calling Repository.Count(1): nil; want not nil")
	}

}

