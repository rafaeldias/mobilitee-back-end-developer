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
	defer db.Close()

	g, err := gorm.Open("mysql", db)
	if err != nil {
		t.Errorf("got error while calling gotm.Open(\"mysql\", %+v): %s, want nil",
			db, err.Error())
	}

	repo := NewRepository(g)

	if !reflect.DeepEqual(repo.db, g) {
		t.Errorf("got db property of Repository: %+v, want: %+v", repo.db, g)
	}
}

func TesRepositoryRead(t *testing.T) {
	testCases := []struct {
	}{}
}
