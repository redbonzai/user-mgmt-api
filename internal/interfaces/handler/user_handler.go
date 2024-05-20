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

// GetUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} user.User
// @Router /users [get]
func (handler *UserHandler) GetUsers(c echo.Context) error {
	users, err := handler.service.GetUsers()
	fmt.Printf("GET USERS : %v+\n", users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} user.User
// @Router /users/{id} [get]
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

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body user.User true "Create User"
// @Success 201 {object} map[string]int
// @Router /users [post]
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

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body user.User true "Update User"
// @Success 200
// @Router /users/{id} [put]
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

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 204
// @Router /users/{id} [delete]
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
