package infrastructure

import (
	echoSwagger "github.com/swaggo/echo-swagger" //nolint:depguard
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/redbonzai/user-management-api/docs"
	"github.com/redbonzai/user-management-api/internal/db"
	"github.com/redbonzai/user-management-api/internal/interfaces/handler"
	"github.com/redbonzai/user-management-api/internal/interfaces/repository"
	internalMiddleware "github.com/redbonzai/user-management-api/internal/middleware"
	"github.com/redbonzai/user-management-api/internal/middleware/authentication"
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

	// Initialize repositories, services, and handlers
	//roleRepo := repository.NewRoleRepository()
	//roleService := services.NewRoleService(roleRepo)
	//roleHandler := handler.NewRoleHandler(roleService)
	//
	//permissionRepo := repository.NewPermissionRepository()
	//permissionService := services.NewPermissionService(permissionRepo)
	//permissionHandler := handler.NewPermissionHandler(permissionService)

	// Public routes
	router.POST("/users/login", userHandler.Login)
	router.POST("/users/register", userHandler.Register)

	// Apply the response interceptor
	router.Use(internalMiddleware.ResponseInterceptor)

	// Protected routes
	protected := router.Group("/v1/users")
	protected.Use(authentication.AuthMiddleware)

	protected.GET("", userHandler.GetUsers)
	protected.GET("/:id", userHandler.GetUser)
	protected.POST("", userHandler.CreateUser)
	protected.PATCH("/:id", userHandler.UpdateUser)
	protected.DELETE("/:id", userHandler.DeleteUser)
	protected.POST("/logout", userHandler.Logout)
	protected.GET("/current-user", userHandler.GetAuthenticatedUser)

	// Role routes
	//protected.GET("/roles", roleHandler.GetRoles)
	//protected.GET("/roles/:id", roleHandler.GetRole)
	//protected.POST("/roles", roleHandler.CreateRole)
	//protected.PUT("/roles/:id", roleHandler.UpdateRole)
	//protected.DELETE("/roles/:id", roleHandler.DeleteRole)
	//protected.POST("/roles/:role_id/users/:user_id", roleHandler.AssignRoleToUser)
	//protected.DELETE("/roles/:role_id/users/:user_id", roleHandler.UnassignRoleFromUser)

	// Permission routes
	//protected.GET("/permissions", permissionHandler.GetPermissions)
	//protected.GET("/permissions/:id", permissionHandler.GetPermission)
	//protected.POST("/permissions", permissionHandler.CreatePermission)
	//protected.PUT("/permissions/:id", permissionHandler.UpdatePermission)
	//protected.DELETE("/permissions/:id", permissionHandler.DeletePermission)
	//protected.POST("/permissions/:permission_id/roles/:role_id", permissionHandler.AssignPermissionToRole)
	//protected.DELETE("/permissions/:permission_id/roles/:role_id", permissionHandler.UnassignPermissionFromRole)

	// Serve Swagger documentation
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	return router
}
