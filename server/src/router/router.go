package router

import (
	"fmt"
	"io"
	"os"
	"log/slog"
	"users-service/src/controller"
	"users-service/src/database"
	"users-service/src/service"
	"users-service/src/config"
	"users-service/src/middleware"

	"github.com/gin-gonic/gin"
)

// Router is a wrapper for the gin.Engine and the address where it is running
type Router struct {
	Engine  *gin.Engine
	Address string
}

// Creates a new router with the configuration provided in the env file
func createRouterFromConfig(cfg *config.Config) *Router {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	router := &Router{
		Engine:  gin.Default(),
		Address: cfg.Host + ":" + cfg.Port,
	}

	router.Engine.Use(middleware.RequestLogger())
	router.Engine.Use(middleware.ErrorHandler())

	return router
}

// Creates a new router with the configuration provided in the env file
func CreateRouter() (*Router, error) {
	cfg := config.LoadConfig()
	r := createRouterFromConfig(cfg)

	user_db, err := database.NewUserMemoryDB()
	if err != nil {
		slog.Error("failed to connect to users database", slog.String("error", err.Error()))
		return nil, err
	}
	interests_db, err := database.NewInterestsMemoryDB()
	if err != nil {
		slog.Error("failed to connect to interests database", slog.String("error", err.Error()))
		return nil, err
	}

	userService := service.CreateUserService(user_db, interests_db)
	userController := controller.CreateUserController(userService)

	r.Engine.POST("/users/register", userController.CreateUser)
	r.Engine.GET("/users/:id", userController.GetUserById)
	r.Engine.GET("/users/register", userController.GetRegisterOptions)

	return r, nil
}

// Runs the router in the address provided in the configuration
func (r *Router) Run() {
	fmt.Println("Running in address: ", r.Address)
	r.Engine.Run(r.Address)
}
