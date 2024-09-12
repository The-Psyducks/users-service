package middleware

import (
	"log/slog"
	"net/http"
	"users-service/src/app_errors"
	"users-service/src/model"

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
				if appErr.Code == http.StatusInternalServerError {
					sendInternalServerErrorResponse(c)
				} else {
					sendErrorResponse(c, appErr.Code, appErr.Message, appErr.Error())
				}
			} else if appValidationError, ok := err.(*app_errors.AppValidationError); ok {
				slog.Error(appValidationError.Message, slog.String("error", appValidationError.Error()))
				sendValidationErrorResponse(c, appValidationError)
			} else {	
				slog.Error("unknown error", slog.String("error", err.Error()))
				sendInternalServerErrorResponse(c)
            }
            c.Abort()
        }
    }
}

func sendValidationErrorResponse(ctx *gin.Context, appValidationError *app_errors.AppValidationError) {
	errorResponse := model.ValidationErrorResponse{
		Type:     "about:blank",
		Title:    appValidationError.Message,
		Status:   appValidationError.Code,
		Detail:   appValidationError.Message,
		Instance: ctx.Request.URL.Path,
		Errors:   appValidationError.Errors,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(appValidationError.Code, errorResponse)
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
