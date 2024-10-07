package mocks

import "github.com/tawanr/ft_matcha/internal/models"

type ProfileModelMock struct{}

func (m *ProfileModelMock) Get(id int) (*models.Profile, error) {
	return nil, nil
}

func (m *ProfileModelMock) Insert(userID int, gender models.GenderType, preferences []models.GenderType, bio string) (*models.Profile, error) {
	return nil, nil
}

func (m *ProfileModelMock) Update(userID int, gender models.GenderType, preferences []models.GenderType, bio string) (*models.Profile, error) {
	return nil, nil
}
