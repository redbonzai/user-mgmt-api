package user

type Repository interface {
	GetAll() ([]User, error)
	GetByID(id int) (User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(id int) (User, error)
}
