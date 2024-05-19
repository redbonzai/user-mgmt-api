package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/redbonzai/user-management-api/internal/domain/user"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service}
}

func (handler *UserHandler) GetUsers(c echo.Context) error {
	users, err := handler.service.GetUsers()
	fmt.Printf("GET USERS : %v+\n", users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

func (handler *UserHandler) GetUser(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid ID")
	}
	retrievedUser, err := handler.service.GetUserByID(id)
	if err != nil {
		return context.JSON(http.StatusNotFound, "User not found")
	}
	return context.JSON(http.StatusOK, retrievedUser)
}

func (handler *UserHandler) CreateUser(context echo.Context) error {
	var createdUser user.User
	if err := context.Bind(&createdUser); err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid input")
	}
	newUser, err := handler.service.CreateUser(createdUser)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.JSON(http.StatusCreated, newUser)
}

func (handler *UserHandler) UpdateUser(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid ID")
	}
	var updatedUser user.User
	if err := context.Bind(&updatedUser); err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid input")
	}
	updatedUser.ID = id
	updated, err := handler.service.UpdateUser(updatedUser)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.JSON(http.StatusOK, updated)
}

func (handler *UserHandler) DeleteUser(context echo.Context) error {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid ID")
	}
	deletedUser, err := handler.service.DeleteUser(id)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err)
	}
	return context.JSON(http.StatusOK, deletedUser)
}
