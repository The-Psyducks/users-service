package router

import (
	"os"
	"io"
	"fmt"
	"log/slog"
	"testing"
	"github.com/gin-gonic/gin"
	"users-service/src/config"
	"users-service/src/controller"
	"users-service/src/middleware"
	"users-service/src/service"
	"users-service/src/database/users_db"
	"users-service/src/database/interests_db"
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

// Creates the databases based on the configuration provided in the env file
func createDatabases(cfg *config.Config) (users_db.UserDatabase, interests_db.InterestsDatabase, error) {
    var user_db users_db.UserDatabase
    var interest_db interests_db.InterestsDatabase
    var err error

    if testing.Testing() {
        user_db, err = users_db.CreateUserMemoryDB()
        if err != nil {
            return nil, nil, fmt.Errorf("failed to create user memory database: %w", err)
        }
        interest_db, err = interests_db.CreateInterestsMemoryDB()
        if err != nil {
            return nil, nil, fmt.Errorf("failed to create interests memory database: %w", err)
        }
    } else {
        user_db, err = users_db.CreateUsersPostgresDB(
            cfg.DatabaseHost,
            cfg.DatabasePort,
            cfg.DatabaseName,
            cfg.DatabasePassword,
            cfg.DatabaseUser)
        if err != nil {
            return nil, nil, fmt.Errorf("failed to connect to users database: %w", err)
        }
        interest_db, err = interests_db.CreateInterestsPostgresDB(
            cfg.DatabaseHost,
            cfg.DatabasePort,
            cfg.DatabaseName,
            cfg.DatabasePassword,
            cfg.DatabaseUser)
        if err != nil {
            return nil, nil, fmt.Errorf("failed to connect to interests database: %w", err)
        }
    }

    return user_db, interest_db, nil
}

// Creates a new router with the configuration provided in the env file
func CreateRouter() (*Router, error) {
    cfg := config.LoadConfig()
    r := createRouterFromConfig(cfg)

    user_db, interest_db, err := createDatabases(cfg)
    if err != nil {
        slog.Error("failed to create databases", slog.String("error", err.Error()))
        return nil, err
    }

    userService := service.CreateUserService(user_db, interest_db)
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
