package models

import (
	"database/sql"
	"errors"
	"slices"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	Age          int
	Name         string
	Gender       GenderType
	PreferMale   bool
	PreferFemale bool
	Bio          string
	Images       []Image
}

type ProfileInterface interface {
	Get(id int) (*Profile, error)
	Insert(userID int, gender GenderType, preferences []int, bio string, age int) (*Profile, error)
	Update(userID int, gender GenderType, preferences []int, bio string, age int) (*Profile, error)
	AddImage(userID int, path Image) error
}

type ProfileModel struct {
	DB *sql.DB
}

func (m *ProfileModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM profile WHERE user_id = ?)`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *ProfileModel) Get(id int) (*Profile, error) {
	var p Profile
	stmt := `SELECT profile.*, users.name from profile join users on profile.user_id = users.id WHERE profile.user_id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&p.UserID, &p.Gender, &p.PreferMale, &p.PreferFemale, &p.Age, &p.Bio, &p.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	caser := cases.Title(language.English)
	p.Name = caser.String(p.Name)

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

func (m *ProfileModel) Insert(userID int, gender GenderType, preferences []int, bio string, age int) (*Profile, error) {
	stmt := `INSERT INTO profile (user_id, gender, prefer_male, prefer_female, bio, age) VALUES (?, ?, ?, ?, ?, ?)
		RETURNING user_id, gender, prefer_male, prefer_female, bio, age`
	preferMale := false
	if slices.Contains(preferences, 0) {
		preferMale = true
	}
	preferFemale := false
	if slices.Contains(preferences, 1) {
		preferFemale = true
	}

	var p Profile
	err := m.DB.QueryRow(stmt, userID, gender, preferMale, preferFemale, bio, age).Scan(&p.UserID, &p.Gender, &p.PreferMale, &p.PreferFemale, &p.Bio, &p.Age)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (m *ProfileModel) Update(userID int, gender GenderType, preferences []int, bio string, age int) (*Profile, error) {
	pf, err := m.Exists(userID)
	if err != nil {
		return nil, err
	}
	if !pf {
		_, err = CreateDefaultProfile(userID, m)
		if err != nil {
			return nil, err
		}
	}
	stmt := `UPDATE profile SET gender = ?, prefer_male = ?, prefer_female = ?, bio = ?, age = ? WHERE user_id = ?
			RETURNING user_id, gender, prefer_male, prefer_female, bio, age`
	preferMale := false
	if slices.Contains(preferences, 0) {
		preferMale = true
	}
	preferFemale := false
	if slices.Contains(preferences, 1) {
		preferFemale = true
	}

	var p Profile

	err = m.DB.QueryRow(stmt, gender, preferMale, preferFemale, bio, age, userID).Scan(&p.UserID, &p.Gender, &p.PreferMale, &p.PreferFemale, &p.Bio, &p.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &p, nil
}

func (m *ProfileModel) AddImage(userID int, path Image) error {
	stmt := `INSERT INTO user_images (user_id, path, priority) VALUES (?, ?, 1)`
	_, err := m.DB.Exec(stmt, userID, path)
	return err
}

func CreateDefaultProfile(userID int, model ProfileInterface) (*Profile, error) {
	profile, err := model.Insert(userID, 0, []int{}, "", -1)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
