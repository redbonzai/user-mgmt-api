package infrastructure

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/redbonzai/user-management-api/docs"
	"github.com/redbonzai/user-management-api/internal/db"
	"github.com/redbonzai/user-management-api/internal/domain/user"
	"github.com/redbonzai/user-management-api/internal/interfaces/handler"
	"github.com/redbonzai/user-management-api/internal/interfaces/repository"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewRouter() *echo.Echo {
	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	userRepo := repository.NewUserRepository(db.DB)
	userService := user.NewService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router.GET("/users", userHandler.GetUsers)
	router.GET("/users/:id", userHandler.GetUser)
	router.POST("/users", userHandler.CreateUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	// Serve Swagger documentation
	router.GET("/swagger/*", echoSwagger.WrapHandler)
	return router

}
