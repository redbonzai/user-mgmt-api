package user

type Repository interface {
	GetAll() ([]User, error)
	GetByID(id int) (User, error)
	Create(user User) (int, error)
	Update(user User) error
	Delete(id int) error
}
