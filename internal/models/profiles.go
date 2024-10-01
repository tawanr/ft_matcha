package models

import "database/sql"

type GenderType int

func (g GenderType) String() string {
	return [...]string{"Male", "Female", "Other", "Not Set"}[g]
}

type Profile struct {
	User        *User
	Gender      string
	Preferences []GenderType
	Bio         string
}

type ProfileModel struct {
	DB *sql.DB
}
