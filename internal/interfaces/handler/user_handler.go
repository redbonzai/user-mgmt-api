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

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.service.GetUsers()
	fmt.Printf("GET USERS : %v+\n", users)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var u user.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}
	id, err := h.service.CreateUser(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	var u user.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}
	u.ID = id
	if err := h.service.UpdateUser(u); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	if err := h.service.DeleteUser(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
