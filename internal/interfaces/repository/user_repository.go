package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/redbonzai/user-management-api/internal/domain/user"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &userRepository{db}
}

func (repository *userRepository) GetAll() ([]user.User, error) {
	var users []user.User
	rows, err := squirrel.Select("id", "name", "email").From("users").RunWith(repository.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var retrievedUser user.User
		if err := rows.Scan(&retrievedUser.ID, &retrievedUser.Name, &retrievedUser.Email); err != nil {
			return nil, err
		}
		users = append(users, retrievedUser)
	}
	return users, nil
}

func (repository *userRepository) GetByID(id int) (user.User, error) {
	var retrievedUser user.User
	err := squirrel.Select("id", "name", "email").From("users").Where(squirrel.Eq{"id": id}).RunWith(repository.db).QueryRow().Scan(&retrievedUser.ID, &retrievedUser.Name, &retrievedUser.Email)
	if err != nil {
		return retrievedUser, err
	}
	return retrievedUser, nil
}

func (repository *userRepository) Create(createdUser user.User) (int, error) {
	var id int
	err := squirrel.Insert("users").Columns("name", "email").Values(createdUser.Name, createdUser.Email).Suffix("RETURNING id").RunWith(repository.db).QueryRow().Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repository *userRepository) Update(updatedUser user.User) error {
	_, err := squirrel.Update("users").Set("name", updatedUser.Name).Set("email", updatedUser.Email).Where(squirrel.Eq{"id": updatedUser.ID}).RunWith(repository.db).Exec()
	return err
}

func (repository *userRepository) Delete(id int) error {
	_, err := squirrel.Delete("users").Where(squirrel.Eq{"id": id}).RunWith(repository.db).Exec()
	return err
}
