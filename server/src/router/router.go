package router

import (
	"io"
	"os"
	"fmt"
	"log/slog"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"users-service/src/config"
	"users-service/src/service"
	"users-service/src/controller"
	"users-service/src/middleware"
	"users-service/src/database/users_db"
	"users-service/src/database/registry_db"
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

// Creates a new database connection using the configuration provided in the env file
func createDBConnection(cfg *config.Config) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	if cfg.Environment == "HEROKU" {
		db, err = sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	} else {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				cfg.DatabaseUser,
				cfg.DatabasePassword,
				cfg.DatabaseHost,
				cfg.DatabasePort,
				cfg.DatabaseName)
	
		db, err = sqlx.Connect("postgres", dsn)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	enableUUIDExtension := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	if _, err := db.Exec(enableUUIDExtension); err != nil {
		return nil, fmt.Errorf("failed to enable uuid extension: %w", err)
	}

	return db, nil
}

// Creates the databases for the users, interests and registry
func createDatabases(cfg *config.Config) (users_db.UserDatabase, registry_db.RegistryDatabase, error) {
	var userDb users_db.UserDatabase
	var registryDb registry_db.RegistryDatabase
	var err error

	db, err := createDBConnection(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	userDb, err = users_db.CreateUsersPostgresDB(db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to users database: %w", err)
	}
	registryDb, err = registry_db.CreateRegistryPostgresDB(db)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to registry database: %w", err)
	}
	return userDb, registryDb, nil
}

// Creates a new router with the configuration provided in the env file
func CreateRouter() (*Router, error) {
	cfg := config.LoadConfig()
	r := createRouterFromConfig(cfg)

	userDb, registryDb, err := createDatabases(cfg)
	if err != nil {
		slog.Error("failed to create databases", slog.String("error", err.Error()))
		return nil, err
	}

	userService := service.CreateUserService(userDb, registryDb)
	userController := controller.CreateUserController(userService)

	public := r.Engine.Group("/")
	{
		public.POST("/users/resolver", userController.ResolveUserEmail)
		
		public.GET("/users/register/locations", userController.GetLocations)
		public.GET("/users/register/interests", userController.GetInterests)
		public.POST("/users/register/:id/send-email", userController.SendVerificationEmail)
		public.POST("/users/register/:id/verify-email", userController.VerifyEmail)
		public.PUT("/users/register/:id/personal-info", userController.AddPersonalInfo)
		public.PUT("/users/register/:id/interests", userController.AddInterests)
		public.POST("/users/register/:id/complete", userController.CompleteRegistry)

		public.POST("/users/login", userController.Login)
	}

	private := r.Engine.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/users/:id", userController.GetUserProfileById)
		private.PUT("/users/profile", userController.ModifyUserProfile)
		
		private.POST("/users/:id/follow", userController.FollowUser)
		private.DELETE("/users/:id/follow", userController.UnfollowUser)
		private.GET("/users/:id/followers", userController.GetFollowers)
		private.GET("/users/:id/following", userController.GetFollowing)
	}

	r.Engine.NoRoute(userController.HandleNoRoute)
	return r, nil
}

// Runs the router in the address provided in the env file
func (r *Router) Run() error {
	fmt.Println("Running in address: ", r.Address)
	return r.Engine.Run(r.Address)
}
