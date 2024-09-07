package middleware

import (
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/model"
	"log/slog"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            if appErr, ok := err.(*app_errors.AppError); ok {
				slog.Error(appErr.Message, slog.String("error", appErr.Error()))
				sendErrorResponse(c, appErr.Code, appErr.Message, appErr.Error())
			} else {	
				slog.Error("unknown error", slog.String("error", err.Error()))
				sendInternalServerErrorResponse(c)
            }
            c.Abort()
        }
    }
}

func sendErrorResponse(ctx *gin.Context, status int, title string, detail string) {
	errorResponse := model.ErrorResponse{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: ctx.Request.URL.Path,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(status, errorResponse)
}

func sendInternalServerErrorResponse(ctx *gin.Context) {
	sendErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", "internal server error - please contact support")
}
