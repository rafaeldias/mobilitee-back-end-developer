package device

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func TestNewRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("got error while calling sqlmock.New(): %s, want nil", err.Error())
	}

	g, err := gorm.Open("mysql", db)
	if err != nil {
		t.Errorf("got error while calling gorm.Open(\"mysql\", %+v): %s, want nil",
			db, err.Error())
	}
	defer g.Close()

	repo := NewRepository(g)

	if !reflect.DeepEqual(repo.db, g) {
		t.Errorf("got db property of Repository: %+v, want: %+v", repo.db, g)
	}
}

func TestRepositoryUpdate(t *testing.T) {
	testCases := []struct {
		id        int
		device    *Device
		wantQuery string
	}{
		{1, &Device{"New Name"}, "UPDATE `devices` SET `name` = \\? WHERE \\(id = \\?\\)"},
		{2, &Device{"New Name"}, "UPDATE `devices` SET `name` = \\? WHERE \\(id = \\?\\)"},
	}

	for _, tc := range testCases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("got error while calling sqlmock.New(): %s, want nil",
				err.Error())
		}

		g, err := gorm.Open("mysql", db)
		if err != nil {
			t.Errorf("got error while calling gotm.Open(\"mysql\", %+v): %s, want nil",
				db, err.Error())
		}

		mock.ExpectBegin()
		mock.ExpectExec(tc.wantQuery).
			WithArgs(tc.device.Name, tc.id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		repo := NewRepository(g)

		err = repo.Update(tc.id, tc.device)
		if err != nil {
			t.Errorf("got error while calling repo.Update(%d, %+v): %s, want nil",
				tc.id, tc.device, err.Error())
		}

		g.Close()
	}
}

func TestRepositoryUpdateError(t *testing.T) {
	testCases := []struct {
		err error
	}{
		{gorm.ErrInvalidSQL},
		{gorm.ErrRecordNotFound},
	}

	for _, tc := range testCases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("got error while creating sqlmock: %s; want nil", err.Error())
		}

		g, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("got error while creating gorm: %s; want nil", err.Error())
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `devices`").WillReturnError(tc.err)
		mock.ExpectRollback()

		repo := NewRepository(g)

		if err := repo.Update(1, &Device{Name: "New Name"}); err == nil {
			t.Error("got error nil while calling repo.Update(0, &Device{}): nil; want not nil")
		}
		g.Close()
	}
}
