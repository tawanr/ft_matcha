package models

import (
	"database/sql"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file:testing.db")
	if err != nil {
		t.Fatal(err)
	}
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		t.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://../../migrations", "testing", driver)
	if err != nil {
		t.Fatal(err)
	}
	m.Up()

	script, err := os.ReadFile("testdata/fixtures.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer db.Close()

		err = os.Remove("testing.db")
		if err != nil {
			t.Fatal(err)
		}
	})
	return db
}
