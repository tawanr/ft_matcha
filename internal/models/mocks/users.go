package mocks

type UserModelMock struct{}

func (m *UserModelMock) Insert(name, email, password string) error {
	return nil
}

func (m *UserModelMock) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModelMock) Exists(id int) (bool, error) {
	return false, nil
}
