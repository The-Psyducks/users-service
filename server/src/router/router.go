package router

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"
	"users-service/src/config"
	"users-service/src/controller"
	"users-service/src/database/registry_db"
	"users-service/src/database/users_db"
	"users-service/src/middleware"
	"users-service/src/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Router is a wrapper for the gin.Engine and the address where it is running
type Router struct {
	Engine  *gin.Engine
	Address string
}

func (r *Router) setNewRelicMiddleware() error {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("users-micro"),
		newrelic.ConfigLicense("ed0d3b23a2f596f67f3c740627feb84aFFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)

	if err != nil {
		return err
	}

	r.Engine.Use(nrgin.Middleware(app))
	return nil
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

	return router
}

// Creates a new database connection using the configuration provided in the env file
func createDBConnection(cfg *config.Config) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	if testing.Testing() {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DatabaseUser,
			cfg.DatabasePassword,
			cfg.DatabaseHost,
			cfg.DatabasePort,
			cfg.DatabaseName)

		db, err = sqlx.Connect("postgres", dsn)
	} else {
		switch cfg.Environment {
		case "HEROKU":
			fallthrough
		case "production":
			db, err = sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
		case "development":
			fallthrough
		case "testing":
			dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				cfg.DatabaseUser,
				cfg.DatabasePassword,
				cfg.DatabaseHost,
				cfg.DatabasePort,
				cfg.DatabaseName)

			db, err = sqlx.Connect("postgres", dsn)
		default:
			return nil, fmt.Errorf("invalid environment: %s", cfg.Environment)
		}
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

	test := false
	if cfg.Environment == "testing" || testing.Testing() {
		test = true
	}

	userDb, err = users_db.CreateUsersPostgresDB(db, test)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to users database: %w", err)
	}
	registryDb, err = registry_db.CreateRegistryPostgresDB(db, test)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to registry database: %w", err)
	}
	return userDb, registryDb, nil
}

func addCorsConfiguration(r *Router) {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	r.Engine.Use(cors.New(config))
}


// Creates a new router with the configuration provided in the env file
func CreateRouter() (*Router, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	r := createRouterFromConfig(cfg)

	userDb, registryDb, err := createDatabases(cfg)
	if err != nil {
		slog.Error("failed to create databases", slog.String("error", err.Error()))
		return nil, err
	}

	if err := r.setNewRelicMiddleware(); err != nil {
		return nil, fmt.Errorf("error setting up newRelic: %w", err)
	}

	r.Engine.Use(middleware.RequestLogger())
	r.Engine.Use(middleware.ErrorHandler())

	addCorsConfiguration(r)

	userService := service.CreateUserService(userDb, registryDb)
	userController := controller.CreateUserController(userService)

	public := r.Engine.Group("/")
	{
		public.POST("/users/resolver", userController.ResolveUserEmail)

		public.GET("/users/info/locations", userController.GetLocations)
		public.GET("/users/info/interests", userController.GetInterests)
		public.POST("/users/register/:id/send-email", userController.SendVerificationEmail)
		public.POST("/users/register/:id/verify-email", userController.VerifyEmail)
		public.PUT("/users/register/:id/personal-info", userController.AddPersonalInfo)
		public.PUT("/users/register/:id/interests", userController.AddInterests)
		public.POST("/users/register/:id/complete", userController.CompleteRegistry)

		public.POST("/users/login", userController.Login)
	}

	private := r.Engine.Group("/")
	private.Use(middleware.AuthMiddleware())
	private.Use(middleware.UserBlockedMiddleware(userService))
	{
		private.GET("/users/:id", userController.GetUserProfileById)
		private.PUT("/users/profile", userController.ModifyUserProfile)
		private.GET("/users/:id/information", userController.GetUserInformation)

		private.POST("/users/:id/follow", userController.FollowUser)
		private.DELETE("/users/:id/follow", userController.UnfollowUser)
		private.GET("/users/:id/followers", userController.GetFollowers)
		private.GET("/users/:id/following", userController.GetFollowing)
		private.POST("/users/:id/block", userController.BlockUser)
		private.POST("/users/:id/unblock", userController.UnblockUser)

		private.GET("/users/search", userController.SearchUsers)

		private.GET("/users/recommendations", userController.RecommendUsers)

		private.GET("/users/all", userController.GetAllUsers)

		private.GET("/users/metrics/registry", userController.GetRegistrationMetrics)
		private.GET("/users/metrics/login", userController.GetLoginMetrics)
		private.GET("/users/metrics/location", userController.GetLocationMetrics)
		private.GET("/users/metrics/blocked", userController.GetUsersBlockedMetrics)
		
	}

	r.Engine.NoRoute(userController.HandleNoRoute)
	return r, nil
}

// Runs the router in the address provided in the env file
func (r *Router) Run() error {
	fmt.Println("Running in address: ", r.Address)
	return r.Engine.Run(r.Address)
}
