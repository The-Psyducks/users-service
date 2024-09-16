package router

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	// "testing"
	"users-service/src/config"
	"users-service/src/controller"
	"users-service/src/database/interests_db"
	"users-service/src/database/users_db"
	"users-service/src/database/registry_db"
	"users-service/src/middleware"
	"users-service/src/service"

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

// Creates the databases based on the configuration provided in the env file
func createDatabases(cfg *config.Config) (users_db.UserDatabase, interests_db.InterestsDatabase, registry_db.RegistryDatabase, error) {
	var userDb users_db.UserDatabase
	var interestDb interests_db.InterestsDatabase
	var registryDb registry_db.RegistryDatabase
	var err error

	// if testing.Testing() || cfg.Environment == "development" {
	// 	userDb, err = users_db.CreateUserMemoryDB()
	// 	if err != nil {
	// 		return nil, nil, nil, fmt.Errorf("failed to create user memory database: %w", err)
	// 	}
	// 	interestDb, err = interests_db.CreateInterestsMemoryDB()
	// 	if err != nil {
	// 		return nil, nil, nil, fmt.Errorf("failed to create interests memory database: %w", err)
	// 	}
	// 	registryDb, err = registry_db.CreateRegistryMemoryDB()
	// 	if err != nil {
	// 		return nil, nil, nil, fmt.Errorf("failed to create registry memory database: %w", err)
	// 	}
	// } else {
	// 	userDb, err = users_db.CreateUsersPostgresDB(
	// 		cfg.DatabaseHost,
	// 		cfg.DatabasePort,
	// 		cfg.DatabaseName,
	// 		cfg.DatabasePassword,
	// 		cfg.DatabaseUser)
	// 	if err != nil {
	// 		return nil, nil, nil, fmt.Errorf("failed to connect to users database: %w", err)
	// 	}
	// 	interestDb, err = interests_db.CreateInterestsPostgresDB(
	// 		cfg.DatabaseHost,
	// 		cfg.DatabasePort,
	// 		cfg.DatabaseName,
	// 		cfg.DatabasePassword,
	// 		cfg.DatabaseUser)
	// 	if err != nil {
	// 		return nil, nil, nil, fmt.Errorf("failed to connect to interests database: %w", err)
	// 	}
	// 	registryDb, err = registry_db.CreateRegistryMemoryDB()
	// 	if err != nil {
	// 		return nil, nil, nil, fmt.Errorf("failed to create registry memory database: %w", err)
	// 	}
	// }
	userDb, err = users_db.CreateUsersPostgresDB(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
		cfg.DatabasePassword,
		cfg.DatabaseUser)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to users database: %w", err)
	}
	interestDb, err = interests_db.CreateInterestsPostgresDB(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
		cfg.DatabasePassword,
		cfg.DatabaseUser)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to interests database: %w", err)
	}
	registryDb, err = registry_db.CreateRegistryPostgresDB(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
		cfg.DatabasePassword,
		cfg.DatabaseUser)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to registry database: %w", err)
	}
	return userDb, interestDb, registryDb, nil
}

// Creates a new router with the configuration provided in the env file
func CreateRouter() (*Router, error) {
	cfg := config.LoadConfig()
	r := createRouterFromConfig(cfg)

	userDb, interestDb, registryDb, err := createDatabases(cfg)
	if err != nil {
		slog.Error("failed to create databases", slog.String("error", err.Error()))
		return nil, err
	}

	userService := service.CreateUserService(userDb, interestDb, registryDb)
	userController := controller.CreateUserController(userService)

	r.Engine.GET("/users/register", userController.GetRegisterOptions)
	
    r.Engine.POST("/users/resolver", userController.ResolveUserEmail)
    r.Engine.POST("/users/register/:id/send-email", userController.SendVerificationEmail)
    r.Engine.POST("/users/register/:id/verify-email", userController.VerifyEmail)
    r.Engine.PUT("/users/register/:id/personal-info", userController.AddPersonalInfo)
    r.Engine.PUT("/users/register/:id/interests", userController.AddInterests)
    r.Engine.POST("/users/register/:id/complete", userController.CompleteRegistry)


	r.Engine.GET("/users/:username", userController.GetUserByUsername)
	r.Engine.POST("/users/login", userController.Login)

	return r, nil
}

// Runs the router in the address provided in the configuration
func (r *Router) Run() error {
	fmt.Println("Running in address: ", r.Address)
	return r.Engine.Run(r.Address)
}
