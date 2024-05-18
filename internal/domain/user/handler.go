package user

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service Service
}

func NewUserHandler(service Service) *UserHandler {
	return &UserHandler{service}
}

// GetUsers retrieves all users
func (handler *UserHandler) GetUsers(context echo.Context) error {
	users, err := handler.service.GetUsers()
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.JSON(http.StatusOK, users)
}

// GetUser retrieves a user by ID
func (handler *UserHandler) GetUser(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid user ID")
	}
	user, err := handler.service.GetUserByID(id)
	if err != nil {
		return context.JSON(http.StatusNotFound, "User not found")
	}
	return context.JSON(http.StatusOK, user)
}

// CreateUser creates a new user
func (handler *UserHandler) CreateUser(context echo.Context) error {
	var createdUser User
	if err := context.Bind(&createdUser); err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid input")
	}
	id, err := handler.service.CreateUser(createdUser)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.JSON(http.StatusCreated, map[string]int{"id": id})
}

// UpdateUser updates an existing user
func (handler *UserHandler) UpdateUser(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid user ID")
	}

	var updatedUser User
	if err := context.Bind(&updatedUser); err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid input")
	}
	updatedUser.ID = id
	if err := handler.service.UpdateUser(updatedUser); err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.NoContent(http.StatusOK)
}

// DeleteUser deletes a user
func (handler *UserHandler) DeleteUser(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid user ID")
	}
	if err := handler.service.DeleteUser(id); err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.NoContent(http.StatusNoContent)
}
