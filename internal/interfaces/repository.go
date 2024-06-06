package interfaces

import "time"

type Repository interface {
	GetAll() ([]User, error)
	GetByUsername(username string) (User, error)
	GetByID(id int) (User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(id int) (User, error)
	GenerateHashFromPassword(password string) (string, error)
	BlacklistToken(token string, expiry time.Time) error
}
