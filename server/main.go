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

	if err := r.Run(); err != nil {
		log.Fatalf("failed to start router: %v", err)
	}
}