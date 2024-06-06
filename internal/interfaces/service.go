package interfaces

import "time"

type Service interface {
	GetUsers() ([]User, error)
	GetUserByUsername(username string) (User, error)
	GetUserByID(id int) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id int) (User, error)
	IsUsernameUnique(username string) (bool, error)
	HashPassword(password string) (string, error)
	Logout(token string, expiry time.Time) error
}
