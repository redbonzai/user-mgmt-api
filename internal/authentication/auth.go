package authentication

import (
	"net/http"
	"os"
	"strings"

	_ "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/redbonzai/user-management-api/pkg/logger"
	"go.uber.org/zap"
)

// JWTMiddleware checks for a valid JWT token
func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(os.Getenv("SECRET_KEY")),
		SigningMethod: "HS256",
		TokenLookup:   "header:Authorization:Bearer",
		ContextKey:    "user",
		ErrorHandler: func(context echo.Context, err error) error {
			logger.Error("Error validating jwt:", zap.Error(err))
			
			return context.JSON(http.StatusUnauthorized, map[string]string{
				"message": "invalid or expired jwt",
			})
		},
	})
}

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		authHeader := context.Request().Header.Get("Authorization")

		if authHeader == "" {
			return context.JSON(http.StatusUnauthorized, "missing or malformed jwt")
		}
		tokenStr := strings.Split(authHeader, " ")[1]
		username, err := ParseToken(tokenStr)
		if err != nil {
			return context.JSON(http.StatusUnauthorized, "invalid or expired jwt")
		}
		context.Set("username", username)
		return next(context)
	}
}

// JWTMiddleware checks for a valid JWT token
//func JWTMiddleware() echo.MiddlewareFunc {
//	return middleware.JWTWithConfig(middleware.JWTConfig{
//		SigningKey:  []byte(os.Getenv("SECRET_KEY")),
//		TokenLookup: "header:Authorization",
//		AuthScheme:  "Bearer",
//	})
//}
//
//// AuthMiddleware checks if the user is authenticated
//func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(context echo.Context) error {
//		authHeader := context.Request().Header.Get("Authorization")
//		if authHeader == "" {
//			return context.JSON(http.StatusUnauthorized, "missing or malformed jwt")
//		}
//		tokenStr := strings.Split(authHeader, " ")[1]
//		username, err := ParseToken(tokenStr)
//		if err != nil {
//			return context.JSON(http.StatusUnauthorized, "invalid or expired jwt")
//		}
//		context.Set("username", username)
//		return next(context)
//	}
//}
