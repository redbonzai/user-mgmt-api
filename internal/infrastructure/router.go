package infrastructure

import (
	echoSwagger "github.com/swaggo/echo-swagger" //nolint:depguard
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/redbonzai/user-management-api/docs"
	"github.com/redbonzai/user-management-api/internal/authentication"
	"github.com/redbonzai/user-management-api/internal/db"
	"github.com/redbonzai/user-management-api/internal/interfaces/handler"
	"github.com/redbonzai/user-management-api/internal/interfaces/repository"
	"github.com/redbonzai/user-management-api/internal/services"
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
	userService := services.NewService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Public routes
	router.POST("/users/login", userHandler.Login)
	router.POST("/users/register", userHandler.Register)

	// Protected routes
	protected := router.Group("/users")
	protected.Use(authentication.AuthMiddleware)

	protected.GET("", userHandler.GetUsers)
	protected.GET("/:id", userHandler.GetUser)
	protected.POST("", userHandler.CreateUser)
	protected.PUT("/:id", userHandler.UpdateUser)
	protected.DELETE("/:id", userHandler.DeleteUser)
	protected.DELETE("/logout", userHandler.Logout)

	// Serve Swagger documentation
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	return router
}
