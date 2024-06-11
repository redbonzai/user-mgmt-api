package authentication

import (
	"fmt"
	"net/http"
	"strings"

	_ "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/redbonzai/user-management-api/internal/db"
	"github.com/redbonzai/user-management-api/pkg/logger"
	"go.uber.org/zap"
)

// JWTMiddleware checks for a valid JWT token
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, "missing or malformed jwt")
			}
			tokenStr := strings.Split(authHeader, " ")[1]

			// Check if token is blacklisted
			if isBlacklisted(tokenStr) {
				return c.JSON(http.StatusUnauthorized, "invalid or expired jwt")
			}

			// Continue with the next handler
			return next(c)
		}
	}
}

func isBlacklisted(token string) bool {
	var count int
	query := `SELECT COUNT(*) FROM token_blacklist WHERE token = $1`
	err := db.DB.QueryRow(query, token).Scan(&count)
	if err != nil {
		logger.Error("Error checking token blacklist: ", zap.Error(err))
		return false
	}
	return count > 0
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
		fmt.Printf("USERNAME: %v\n", username)

		if err != nil {
			logger.Error("Error parsing token: ", zap.Error(err))
			return context.JSON(http.StatusUnauthorized, "invalid or expired jwt")
		}
		context.Set("username", username)
		return next(context)
	}
}
