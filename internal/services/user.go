package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/redbonzai/user-management-api/internal/interfaces"
)

type service struct {
	repo interfaces.Repository
}

func NewService(repo interfaces.Repository) interfaces.Service {
	return &service{repo}
}

func (service *service) GetUsers() ([]interfaces.User, error) {

	return service.repo.GetAll()
}

func (service *service) GetUserByUsername(username string) (interfaces.User, error) {
	return service.repo.GetByUsername(username)
}

func (service *service) GetUserByID(id int) (interfaces.User, error) {
	return service.repo.GetByID(id)
}

func (service *service) CreateUser(user interfaces.User) (interfaces.User, error) {
	return service.repo.Create(user)
}

func (service *service) UpdateUser(user interfaces.User) (interfaces.User, error) {
	return service.repo.Update(user)
}

func (service *service) DeleteUser(id int) (interfaces.User, error) {
	return service.repo.Delete(id)
}

func (service *service) IsUsernameUnique(username string) (bool, error) {
	_, err := service.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil // Username is unique
		}
		return false, err // Some other error occurred
	}
	return false, nil // Username already exists
}

func (service *service) HashPassword(password string) (string, error) {
	return service.repo.GenerateHashFromPassword(password)
}

func (service *service) Logout(token string, expiry time.Time) error {
	return service.repo.BlacklistToken(token, expiry)
}
