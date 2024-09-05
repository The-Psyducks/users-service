package main

import (
	"log"
	"users-service/src/router"
)

func main () {
	r, err := router.CreateRouter()

	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	r.Run()
}