package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	// _ "github.com/tursodatabase/go-libsql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, password, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return ErrDuplicateEmail
		}
		return nil
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, password FROM users WHERE email = ?`
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, ErrInvalidCredentials
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
