package middleware

import (
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserBlockedMiddleware(service *service.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool("session_user_admin") {
			c.Next()
			return
		}

		sessionUserId := c.GetString("session_user_id")
		userId, _ := uuid.Parse(sessionUserId)

		isBlocked, err := service.CheckIfUserIsBlocked(userId)
		if err != nil {
			err := app_errors.NewAppError(http.StatusInternalServerError, "Error checking if user is blocked", err)
			_ = c.AbortWithError(err.Code,err)
			c.Next()
			return
		}

		if isBlocked {
			err := app_errors.NewAppError(http.StatusUnauthorized, "User is blocked", nil)
			_ = c.AbortWithError(err.Code,err)
			c.Next()
			return
		}

		c.Next()
	}
}