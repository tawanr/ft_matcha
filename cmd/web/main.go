package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/tursodatabase/go-libsql"
	"github.org/tawanr/ft_matcha/internal/models"
)

type application struct {
	logger        *slog.Logger
	db            *sql.DB
	users         *models.UserModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "file:data.db", "Data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		db:            db,
		users:         &models.UserModel{DB: db},
		templateCache: templateCache,
	}

	app.logger.Info("Starting server", slog.String("addr", *addr))
	logger.Error(http.ListenAndServe(*addr, app.routes()).Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
