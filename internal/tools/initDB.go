package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"

	_ "github.com/tursodatabase/go-libsql"
)

func run(path string) (err error) {
	db, err := sql.Open("libsql", "file:"+path+"/data.db")
	if err != nil {
		return err
	}
	defer func() {
		if closeError := db.Close(); closeError != nil {
			slog.Error("Error closing database", closeError)
			err = closeError
		}
	}()
	if err != nil {
		return err
	}

	// _, err = db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return err
	}
	for rows.Next() {
		fmt.Println(rows)
	}
	return nil
}

func main() {
	path := flag.String("path", "data/libsql", "Path to locate database file")
	slog.Info("Initializing database at " + *path + "/data.db")
	if err := run(*path); err != nil {
		panic(err)
	}
}
