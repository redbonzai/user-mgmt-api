package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	//nolint:depguard
	//nolint:depguard
	"github.com/redbonzai/user-management-api/internal/authentication"
	"github.com/redbonzai/user-management-api/internal/interfaces"
	"github.com/redbonzai/user-management-api/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	service interfaces.Service
}

func NewUserHandler(service interfaces.Service) *UserHandler {
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
func (handler *UserHandler) GetUsers(context echo.Context) error {
	users, err := handler.service.GetUsers()
	fmt.Printf("GET USERS : %v+\n", users)
	if err != nil {
		logger.Error("Error retrieving users: ", zap.Error(err))
		return context.JSON(http.StatusInternalServerError, err)
	}
	logger.Info("Users retrieved", zap.Int("count", len(users)))
	return context.JSON(http.StatusOK, users)
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
		logger.Error("Invalid User ID: ", zap.Error(err))
		return context.JSON(http.StatusBadRequest, "Invalid ID")
	}
	retrievedUser, err := handler.service.GetUserByID(id)
	if err != nil {
		logger.Error("Error retrieving user: ", zap.Error(err))
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
	var createdUser interfaces.User
	if err := context.Bind(&createdUser); err != nil {
		logger.Error("Invalid input: ", zap.Error(err))
		return context.JSON(http.StatusBadRequest, "Invalid input")
	}
	newUser, err := handler.service.CreateUser(createdUser)
	if err != nil {
		logger.Error("Error creating user: ", zap.Error(err))
		return context.JSON(http.StatusInternalServerError, err)
	}
	logger.Info("User created", zap.Int("userID", newUser.ID), zap.String("name", newUser.Name))
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
		logger.Error("Invalid User ID: ", zap.Error(err))
		return context.JSON(http.StatusBadRequest, "Invalid ID")
	}

	var updatedUser interfaces.User
	if err := context.Bind(&updatedUser); err != nil {
		logger.Error("Invalid input: ", zap.Error(err))
		return context.JSON(http.StatusBadRequest, "Invalid input")
	}
	updatedUser.ID = id
	updated, err := handler.service.UpdateUser(updatedUser)
	if err != nil {
		logger.Error("Error updating user: ", zap.Error(err))
		return context.JSON(http.StatusInternalServerError, err)
	}
	logger.Info("User updated", zap.Int("userID", updatedUser.ID), zap.String("name", updatedUser.Name))
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
		logger.Error("Invalid user ID", zap.Error(err))
		return context.JSON(http.StatusBadRequest, "Invalid ID")
	}
	deletedUser, err := handler.service.DeleteUser(id)
	if err != nil {
		logger.Error("Error deleting user", zap.Int("userID", id), zap.Error(err))
		return context.JSON(http.StatusInternalServerError, err)
	}
	logger.Info("User deleted", zap.Int("userID", deletedUser.ID), zap.String("name", deletedUser.Name))
	return context.JSON(http.StatusOK, deletedUser)
}

// Login godoc
// @Summary Login a user
// @Description Login a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body user.LoginRequest true "Credentials"
// @Success 200 {object} map[string]string
// @Router /v1/login [post]
func (handler *UserHandler) Login(context echo.Context) error {
	var loginRequest interfaces.LoginRequest
	if err := context.Bind(&loginRequest); err != nil {

		logger.Error("Invalid input: ", zap.Error(err))
		return context.JSON(http.StatusBadRequest, "Invalid input")
	}

	user, err := handler.service.GetUserByUsername(loginRequest.Username)
	logger.Info("User retrieved", zap.Any("user", user))

	if err != nil || user.Password == "" || user.Username == "" {
		logger.Error(
			"Invalid username or password: ",
			zap.Error(err),
			zap.String("username", loginRequest.Username),
			zap.String("password", loginRequest.Password),
		)
		return context.JSON(http.StatusUnauthorized, "Invalid username or password")
	}

	// Compare the hashed password with the login password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return context.JSON(http.StatusUnauthorized, "Invalid password")
	}

	token, err := authentication.GenerateToken(user.Username)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Failed to generate token")
	}

	return context.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

// Logout godoc
// @Summary Logout a user
// @Description Logout a user and blacklist the token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body user.LoginRequest true "Credentials"
// @Success 200 {object} map[string]string
// @Router /v1/logout [post]
func (handler *UserHandler) Logout(context echo.Context) error {
	userToken := context.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	exp := int64(claims["exp"].(float64))

	// Blacklist the token
	err := handler.service.Logout(userToken.Raw, time.Unix(exp, 0))
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Failed to logout")
	}

	return context.JSON(http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body user.User true "User"
// @Success 201 {object} user.User
// @Router /v1/register [post]
func (handler *UserHandler) Register(context echo.Context) error {
	var registerRequest interfaces.RegisterRequest
	if err := context.Bind(&registerRequest); err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid input")
	}

	// Check for unique username
	isUnique, err := handler.service.IsUsernameUnique(registerRequest.Username)
	if err != nil || !isUnique {
		return context.JSON(http.StatusConflict, "Username already exists")
	}

	// Hash the password using the userRepository method
	hashedPassword, err := handler.service.HashPassword(registerRequest.Password)
	logger.Info("Hashed password", zap.String("password", hashedPassword))

	if err != nil {
		logger.Error("Failed to hash password: ", zap.Error(err), zap.String("password", registerRequest.Password))
		return context.JSON(http.StatusInternalServerError, "Failed to hash password")
	}

	// Create the user
	newUser := interfaces.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Username: registerRequest.Username,
		Password: hashedPassword,
		Status:   registerRequest.Status,
	}
	//*newUser.Status = "active"

	createdUser, err := handler.service.CreateUser(newUser)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Failed to create user")
	}

	return context.JSON(http.StatusCreated, createdUser)
}

func (handler *UserHandler) GetAuthenticatedUser(context echo.Context) error {
	userToken := context.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	user, err := handler.service.GetUserByUsername(username)
	if err != nil {
		return context.JSON(http.StatusNotFound, "User not found")
	}

	return context.JSON(http.StatusOK, user)
}
