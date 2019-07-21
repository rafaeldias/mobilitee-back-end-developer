package device

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func setupDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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
	db, _ := setupDB(t)
	defer db.Close()

	repo := NewRepository(db)

	if !reflect.DeepEqual(repo.db, db) {
		t.Errorf("got db property of Repository: %+v, want: %+v", repo.db, db)
	}
}

func TestRepositoryWrite(t *testing.T) {
	testCases := []struct {
		device    *Device
		wantQuery string
		wantID    int
	}{
		{
			&Device{Name: "New Name", Model: "Android", User: 1, CreatedAt: time.Now()},
			"INSERT INTO `devices`",
			1,
		},
	}

	for _, tc := range testCases {
		db, mock := setupDB(t)

		mock.ExpectBegin()
		mock.ExpectExec(tc.wantQuery).
			WithArgs(tc.device.Name, tc.device.Model, tc.device.User, tc.device.CreatedAt).
			WillReturnResult(sqlmock.NewResult(int64(tc.wantID), 1))
		mock.ExpectCommit()

		repo := NewRepository(db)

		id, err := repo.Write(tc.device)
		if err != nil {
			t.Fatalf("got error while calling Repository.Write(%+v): %s, want nil",
				tc.device, err.Error())
		}

		if id != tc.wantID {
			t.Errorf("got id from calling Repository.Write(%+v): %d, want %d", tc.device,
				id, tc.wantID)
		}

		db.Close()
	}
}

func TestRepositoryWriteError(t *testing.T) {

		db, mock := setupDB(t)
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `devices`").WillReturnError(gorm.ErrInvalidSQL)
		mock.ExpectRollback()

		repo := NewRepository(db)

		if _, err := repo.Write(&Device{}); err == nil {
			t.Error("got error nil while calling Repository.Write(&Device{}); want not nil")
		}

}

func TestRepositoryExchange(t *testing.T) {
	testCases := []struct {
		device		*Device
		old		*Device
		wantQueryWrite  string
		wantQueryUpdate string
		wantID		int
		wantOldID	int
	}{
		{
			&Device{Name: "New Name", Model: "Android", User: 1, CreatedAt: time.Now()},
			&Device{ID: 2},
			"INSERT INTO `devices`",
			"UPDATE `devices`",
			1,
			2,
		},
	}

	for _, tc := range testCases {
		db, mock := setupDB(t)

		mock.ExpectBegin()
		mock.ExpectExec(tc.wantQueryWrite).
			WithArgs(tc.device.Name, tc.device.Model, tc.device.User, tc.device.CreatedAt).
			WillReturnResult(sqlmock.NewResult(int64(tc.wantID), 1))
		mock.ExpectExec(tc.wantQueryUpdate).
			WithArgs(tc.wantID, tc.wantOldID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := NewRepository(db)

		id, err := repo.Exchange(tc.old, tc.device)
		if err != nil {
			t.Fatalf("got error while calling Repository.Exchange(%+v, %+v): %s, want nil",
				tc.old, tc.device, err.Error())
		}

		if id != tc.wantID {
			t.Errorf("got id from calling Repository.Exchange(%+v, %+v): %d, want %d",
				tc.old, tc.device, id, tc.wantID)
		}

		db.Close()
	}
}

type transactionErrorMock func(sqlmock.Sqlmock)

func TestRepositoryExchangeError(t *testing.T) {
		nw := &Device{Name: "New Name", Model: "Android", User: 1, CreatedAt: time.Now()}
		old := &Device{ID: 3}
		testCases := []transactionErrorMock{
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("Internal Transaction error"))
			},
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `devices`").
					WithArgs(nw.Name, nw.Model, nw.User, nw.CreatedAt).
					WillReturnError(gorm.ErrInvalidSQL)
				mock.ExpectRollback()
			},
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `devices`").
					WithArgs(nw.Name, nw.Model, nw.User, nw.CreatedAt). 
					WillReturnResult(sqlmock.NewResult(int64(nw.ID), 1))
				mock.ExpectExec("UPDATE `devices`").
					WithArgs(nw.ID, old.ID).
					WillReturnError(gorm.ErrInvalidSQL)
				mock.ExpectRollback()
			},
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `devices`").
					WithArgs(nw.Name, nw.Model, nw.User, nw.CreatedAt). 
					WillReturnResult(sqlmock.NewResult(int64(nw.ID), 1))
				mock.ExpectExec("UPDATE `devices`").
					WithArgs(nw.ID, old.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(gorm.ErrInvalidSQL)
			},
		}

		for _, tc := range testCases {
			db, mock := setupDB(t)

			tc(mock)

			repo := NewRepository(db)

			if _, err := repo.Exchange(old, nw); err == nil {
				t.Errorf("got error nil while calling Repository.Exchange(%+v, %+v); want not nil",
					old, nw)
			}

			db.Close()
		}

}
