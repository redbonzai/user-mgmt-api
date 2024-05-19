package repository

import (
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/redbonzai/user-management-api/internal/domain/user"
	"github.com/redbonzai/user-management-api/pkg/logger"
	"go.uber.org/zap"
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
		logger.Error("Error building user query:", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var retrievedUser user.User
		if err := rows.Scan(&retrievedUser.ID, &retrievedUser.Name, &retrievedUser.Email); err != nil {
			logger.Error("Error scanning user row:", zap.Error(err))
			return nil, err
		}
		users = append(users, retrievedUser)
	}
	return users, nil
}

func (repository *userRepository) GetByID(id int) (user.User, error) {
	var retrievedUser user.User
	query, args, err := squirrel.Select("id", "name", "email").
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		logger.Error("Error building SQL query:", zap.Error(err))
		return retrievedUser, err
	}

	err = repository.db.QueryRow(query, args...).Scan(&retrievedUser.ID, &retrievedUser.Name, &retrievedUser.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("User not found", zap.Int("userID", id))
			return retrievedUser, err
		}
		logger.Error("Error retrieving user with ID:", zap.Int("userID", id), zap.Error(err))
		return retrievedUser, err
	}
	return retrievedUser, nil
}

func (repository *userRepository) Create(createdUser user.User) (user.User, error) {
	query, args, err := squirrel.Insert("users").
		Columns("name", "email").
		Values(createdUser.Name, createdUser.Email).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar). // Ensure PostgreSQL-compatible placeholders
		ToSql()
	if err != nil {
		logger.Error("Error building SQL query:", zap.Error(err))
		return createdUser, err
	}

	err = repository.db.QueryRow(query, args...).Scan(&createdUser.ID)
	if err != nil {
		logger.Error("Error creating user:", zap.Error(err))
		return createdUser, err
	}

	return createdUser, nil
}

func (repository *userRepository) Update(updatedUser user.User) (user.User, error) {
	query, args, err := squirrel.Update("users").
		Set("name", updatedUser.Name).
		Set("email", updatedUser.Email).
		Where(squirrel.Eq{"id": updatedUser.ID}).
		PlaceholderFormat(squirrel.Dollar). // Ensure PostgreSQL-compatible placeholders
		ToSql()
	if err != nil {
		logger.Error("Error building SQL query: %v", zap.Error(err))
		return updatedUser, err
	}

	_, err = repository.db.Exec(query, args...)
	if err != nil {
		logger.Error("Error updating user:", zap.Error(err))
		return updatedUser, err
	}

	return repository.GetByID(updatedUser.ID)
}

func (repository *userRepository) Delete(id int) (user.User, error) {
	deletedUser, err := repository.GetByID(id)
	if err != nil {
		logger.Error("Error retrieving user to delete:", zap.Error(err))
		return deletedUser, err
	}

	logger.Info("Retrieved Deleted user: ", zap.Any("deletedUser", deletedUser))

	query, args, err := squirrel.Delete("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar). // Ensure PostgreSQL-compatible placeholders
		ToSql()
	if err != nil {
		logger.Error("Error building SQL query: ", zap.Error(err))
		return deletedUser, err
	}

	_, err = repository.db.Exec(query, args...)
	if err != nil {
		logger.Error("Error deleting user:", zap.Error(err))
		return deletedUser, err
	}

	return deletedUser, nil
}
