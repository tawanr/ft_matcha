package models

import (
	"database/sql"
	"errors"
	"slices"
)

type GenderType int
type Gender int
type Image string

func (g Gender) String() string {
	return [...]string{"Male", "Female"}[g]
}

func (g GenderType) String() string {
	return [...]string{"Male", "Female", "Other", "Not Set"}[g]
}

type Profile struct {
	UserID       int64
	Gender       GenderType
	PreferMale   bool
	PreferFemale bool
	Bio          string
	Images       []Image
}

type ProfileInterface interface {
	Get(id int) (*Profile, error)
	Insert(userID int, gender GenderType, preferences []GenderType, bio string) (*Profile, error)
	Update(userID int, gender GenderType, preferences []GenderType, bio string) (*Profile, error)
}

type ProfileModel struct {
	DB *sql.DB
}

func (m *ProfileModel) Get(id int) (*Profile, error) {
	var p Profile
	stmt := `SELECT * FROM profile WHERE user_id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&p.UserID, &p.Gender, &p.PreferMale, &p.PreferFemale, &p.Bio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	stmt = `SELECT path FROM user_images WHERE user_id = ? ORDER BY priority ASC`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &p, nil
		}
		return nil, err
	}
	for rows.Next() {
		var path Image
		err := rows.Scan(&path)
		if err != nil {
			return nil, err
		}
		p.Images = append(p.Images, path)
	}
	return &p, nil
}

func (m *ProfileModel) Insert(userID int, gender GenderType, preferences []GenderType, bio string) (*Profile, error) {
	stmt := `INSERT INTO profile (user_id, gender, prefer_male, prefer_female, bio) VALUES (?, ?, ?, ?, ?)
		RETURNING user_id, gender, prefer_male, prefer_female, bio`
	preferMale := false
	if slices.Contains(preferences, 0) {
		preferMale = true
	}
	preferFemale := false
	if slices.Contains(preferences, 1) {
		preferFemale = true
	}

	var p Profile
	err := m.DB.QueryRow(stmt, userID, gender, preferMale, preferFemale, bio).Scan(&p.UserID, &p.Gender, &p.PreferMale, &p.PreferFemale, &p.Bio)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (m *ProfileModel) Update(userID int, gender GenderType, preferences []GenderType, bio string) (*Profile, error) {
	stmt := `UPDATE profile SET gender = ?, prefer_male = ?, prefer_female = ?, bio = ? WHERE user_id = ?
			RETURNING user_id, gender, prefer_male, prefer_female, bio`
	preferMale := false
	if slices.Contains(preferences, 0) {
		preferMale = true
	}
	preferFemale := false
	if slices.Contains(preferences, 1) {
		preferFemale = true
	}

	var p Profile

	err := m.DB.QueryRow(stmt, gender, preferMale, preferFemale, bio, userID).Scan(&p.UserID, &p.Gender, &p.PreferMale, &p.PreferFemale, &p.Bio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &p, nil
}
