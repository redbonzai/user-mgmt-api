package user

type Service interface {
	GetUsers() ([]User, error)
	GetUserByID(id int) (User, error)
	CreateUser(user User) (int, error)
	UpdateUser(user User) error
	DeleteUser(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetUsers() ([]User, error) {
	return s.repo.GetAll()
}

func (s *service) GetUserByID(id int) (User, error) {
	return s.repo.GetByID(id)
}

func (s *service) CreateUser(user User) (int, error) {
	return s.repo.Create(user)
}

func (s *service) UpdateUser(user User) error {
	return s.repo.Update(user)
}

func (s *service) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
