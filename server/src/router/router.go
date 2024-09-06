package router

import (
	"fmt"
	"users-service/src/controller"
	"users-service/src/database"
	"users-service/src/service"

	"github.com/gin-gonic/gin"
)

// Router is a wrapper for the gin.Engine and the address where it is running
type Router struct {
	Engine  *gin.Engine
	Address string
}

// Creates a new router with the configuration provided in the env file
func CreateRouter() (*Router, error) {
	r := Router{
		Engine:  gin.Default(),
		Address: "0.0.0.0:8080",
	}

	user_db, err := database.NewUserMemoryDB()
	if err != nil {
		return nil, err
	}
	interests_db, err := database.NewInterestsMemoryDB()
	if err != nil {
		return nil, err
	}

	userService := service.CreateUserService(user_db, interests_db)
	userController := controller.CreateUserController(userService)

	r.Engine.POST("/users/register", userController.CreateUser)
	r.Engine.GET("/users/:id", userController.GetUserById)
	r.Engine.GET("/users/register", userController.GetRegisterOptions)

	return &r, nil
}

// Runs the router in the address provided in the configuration
func (r *Router) Run() {
	fmt.Println("Running in address: ", r.Address)
	r.Engine.Run(r.Address)
}
