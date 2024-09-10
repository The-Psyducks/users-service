package router

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"users-service/src/config"
	"users-service/src/controller"
	"users-service/src/middleware"
	"users-service/src/service"

	"users-service/src/database/users_db"
	"users-service/src/database/interests_db"

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

	user_db, err := users_db.NewUserMemoryDB()
	if err != nil {
		slog.Error("failed to connect to users database", slog.String("error", err.Error()))
		return nil, err
	}
	interests_db, err := interests_db.NewInterestsMemoryDB()
	if err != nil {
		slog.Error("failed to connect to interests database", slog.String("error", err.Error()))
		return nil, err
	}

	userService := service.CreateUserService(user_db, interests_db)
	userController := controller.CreateUserController(userService)

	r.Engine.GET("/users/register", userController.GetRegisterOptions)
	r.Engine.POST("/users/register", userController.CreateUser)
	
	r.Engine.GET("/users/:username", userController.GetUserByUsername)

	r.Engine.POST("/users/login", userController.Login)

	return r, nil
}

// Runs the router in the address provided in the configuration
func (r *Router) Run() error {
	fmt.Println("Running in address: ", r.Address)
	return r.Engine.Run(r.Address)
}
