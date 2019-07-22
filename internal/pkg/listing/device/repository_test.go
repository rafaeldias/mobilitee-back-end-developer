package device

import (
	"errors"
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
		t.Errorf("got error while calling gotm.Open(\"mysql\", %+v): %s, want nil",
			db, err.Error())
	}
	defer g.Close()

	repo := NewRepository(g)

	if !reflect.DeepEqual(repo.db, g) {
		t.Errorf("got db property of Repository: %+v, want: %+v", repo.db, g)
	}
}

func TestRepositoryRead(t *testing.T) {
	testCases := []struct {
		ID          int
		wantQuery   string
		wantDevices []*Device
	}{
		{0, "SELECT \\* FROM `devices` WHERE `devices`.`deleted_at` IS NULL", []*Device{&Device{ID: 0}}},
		{1, "SELECT \\* FROM `devices` WHERE `devices`.`deleted_at` IS NULL AND \\(\\(`devices`\\.`id` = 1\\)\\)", []*Device{&Device{ID: 1}}},
		{2, "SELECT \\* FROM `devices` WHERE `devices`.`deleted_at` IS NULL AND \\(\\(`devices`\\.`id` = 2\\)\\)", []*Device{&Device{ID: 2}}},
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

		rows := sqlmock.NewRows([]string{"id"}).AddRow(tc.ID)
		mock.ExpectQuery(tc.wantQuery).WillReturnRows(rows)

		repo := NewRepository(g)

		devices, err := repo.Read(tc.ID)
		if err != nil {
			t.Errorf("got error while calling repo.Read(%d): %s, want nil",
				tc.ID, err.Error())
		}

		if !reflect.DeepEqual(tc.wantDevices, devices) {
			t.Errorf("got devices from repo.Read(%d): %+v, want: %+v", tc.ID,
				devices, tc.wantDevices)

		}

		g.Close()
	}
}

func TestRepositoryReadError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("got error while creating sqlmock: %s; want nil", err.Error())
	}

	g, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("got error while creating gorm: %s; want nil", err.Error())
	}

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("Testing Repository Error"))

	repo := NewRepository(g)

	_, err = repo.Read(0)
	if err == nil {
		t.Error("got error nil while calling repo.Read(0): nil; want not nil")
	}
}
