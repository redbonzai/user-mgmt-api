package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/redbonzai/user-management-api/internal/interfaces"
	"github.com/redbonzai/user-management-api/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) interfaces.Repository {
	return &userRepository{db}
}

func (repository *userRepository) GetAll() ([]interfaces.User, error) {
	var users []interfaces.User
	rows, err := squirrel.
		Select("id", "name", "email", "status", "username", "password").
		From("users").
		OrderBy("id").
		RunWith(repository.db).
		Query()
	if err != nil {
		logger.Error("Error building user query:", zap.Error(err))
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error("Error closing rows:", zap.Error(err))
		}
	}(rows)

	for rows.Next() {
		var retrievedUser interfaces.User
		if err := rows.Scan(
			&retrievedUser.ID,
			&retrievedUser.Name,
			&retrievedUser.Email,
			&retrievedUser.Status,
			&retrievedUser.Username,
			&retrievedUser.Password,
		); err != nil {
			logger.Error("Error scanning user row:", zap.Error(err))
			return nil, err
		}
		users = append(users, retrievedUser)
	}
	return users, nil
}

func (repository *userRepository) GetByUsername(username string) (interfaces.User, error) {
	var retrievedUser interfaces.User
	query, args, err := squirrel.
		Select("id", "name", "email", "status", "username", "password").
		From("users").
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		logger.Error("Error building SQL query:", zap.Error(err))
		return retrievedUser, err
	}

	err = repository.
		db.QueryRow(query, args...).
		Scan(
			&retrievedUser.ID,
			&retrievedUser.Name,
			&retrievedUser.Email,
			&retrievedUser.Status,
			&retrievedUser.Username,
			&retrievedUser.Password,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("User not found", zap.String("username", username))
			return retrievedUser, err
		}
		logger.Error("Error retrieving user with username:", zap.String("username", username), zap.Error(err))
		return retrievedUser, err
	}
	return retrievedUser, nil
}

func (repository *userRepository) GetByID(id int) (interfaces.User, error) {
	var retrievedUser interfaces.User
	query, args, err := squirrel.
		Select("id", "name", "email", "status", "username", "password").
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		logger.Error("Error building SQL query:", zap.Error(err))
		return retrievedUser, err
	}

	err = repository.
		db.QueryRow(query, args...).
		Scan(
			&retrievedUser.ID,
			&retrievedUser.Name,
			&retrievedUser.Email,
			&retrievedUser.Status,
			&retrievedUser.Username,
			&retrievedUser.Password,
		)
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

func (repository *userRepository) Create(createdUser interfaces.User) (interfaces.User, error) {
	query, args, err := squirrel.Insert("users").
		Columns("name", "email", "status", "username", "password").
		Values(createdUser.Name, createdUser.Email, createdUser.Status, createdUser.Username, createdUser.Password).
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

func (repository *userRepository) Update(updatedUser interfaces.User) (interfaces.User, error) {
	query, args, err := squirrel.Update("users").
		Set("name", updatedUser.Name).
		Set("email", updatedUser.Email).
		Set("status", updatedUser.Status).
		Set("status", updatedUser.Username).
		Set("status", updatedUser.Password).
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

func (repository *userRepository) Delete(id int) (interfaces.User, error) {
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

func (repository *userRepository) GenerateHashFromPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// BlacklistToken blacklists a given token until its expiration
func (repository *userRepository) BlacklistToken(token string, expiry time.Time) error {
	_, err := repository.db.Exec("INSERT INTO token_blacklist (token, expiry) VALUES ($1, $2)", token, expiry)
	return err
}
