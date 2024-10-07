package models

import (
	"database/sql"
)

type Models struct {
	Users    UserInterface
	Profiles ProfileInterface
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Users:    &UserModel{DB: db},
		Profiles: &ProfileModel{DB: db},
	}
}
