package device

import (
	"time"
	"testing"
	"reflect"

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

func TestRepositoryRemove(t *testing.T) {
	testCases := []struct {
		device    *Device
		wantQuery string
	}{
		{
			&Device{ID: 1, DeletedAt: time.Now()},
			"UPDATE `devices`",
		},
	}

	for _, tc := range testCases {
		db, mock := setupDB(t)

		mock.ExpectBegin()
		mock.ExpectExec(tc.wantQuery).
			WithArgs(tc.device.DeletedAt, tc.device.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		repo := NewRepository(db)

		if err := repo.Remove(tc.device); err != nil {
			t.Fatalf("got error while calling Repository.Remove(%+v): %s, want nil",
				tc.device, err.Error())
		}

		db.Close()
	}
}

func TestRepositoryRemoveError(t *testing.T) {
		db, mock := setupDB(t)
		defer db.Close()

		d := &Device{ID: 1, DeletedAt: time.Now()}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `devices`").
			WithArgs(d.DeletedAt, d.ID).
			WillReturnError(gorm.ErrInvalidSQL)
		mock.ExpectRollback()

		repo := NewRepository(db)

		if err := repo.Remove(d); err == nil {
			t.Errorf("got error nil while calling Repository.Remove(%+v); want not nil", d)
		}
}




