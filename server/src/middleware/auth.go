package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"users-service/src/app_errors"
	"users-service/src/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Error("Authorization header is required")
			err := app_errors.NewAppError(http.StatusUnauthorized, "Authorization header is required", fmt.Errorf("authorization header is required"))
			_ = c.AbortWithError(err.Code,err)
			c.Next()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			slog.Error("Invalid authorization header")
			err := app_errors.NewAppError(http.StatusUnauthorized, "Invalid authorization header", fmt.Errorf("invalid authorization header"))
			_ = c.AbortWithError(err.Code, err)
			c.Next()
			return
		}

		tokenString := bearerToken[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			slog.Error("Invalid token")
			err := app_errors.NewAppError(http.StatusUnauthorized, "Invalid token", err)
			_ = c.AbortWithError(err.Code, err)
			return
		}

		c.Set("session_user_id", claims.UserId)

		c.Next()
	}
}