package user

type Service interface {
	GetUsers() ([]User, error)
	GetUserByID(id int) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id int) (User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (service *service) GetUsers() ([]User, error) {
	return service.repo.GetAll()
}

func (service *service) GetUserByID(id int) (User, error) {
	return service.repo.GetByID(id)
}

func (service *service) CreateUser(user User) (User, error) {
	return service.repo.Create(user)
}

func (service *service) UpdateUser(user User) (User, error) {
	return service.repo.Update(user)
}

func (service *service) DeleteUser(id int) (User, error) {
	return service.repo.Delete(id)
}
