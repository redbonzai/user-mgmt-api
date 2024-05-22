package handler

import (
	"github.com/redbonzai/user-management-api/internal/domain/user"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUsers() ([]user.User, error) {
	args := m.Called()
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id int) (user.User, error) {
	args := m.Called(id)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserService) CreateUser(u user.User) (user.User, error) {
	args := m.Called(u)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(u user.User) (user.User, error) {
	args := m.Called(u)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(id int) (user.User, error) {
	args := m.Called(id)
	return args.Get(0).(user.User), args.Error(1)
}
