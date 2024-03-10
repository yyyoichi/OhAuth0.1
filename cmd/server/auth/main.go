package main

import (
	"context"
	"log"

	"github.com/yyyoichi/OhAuth0.1/internal/auth"
)

func main() {
	ctx := context.Background()
	service, err := auth.NewService(ctx, auth.Config{
		DatabaseServerURL: "http://localhost:3306",
	})
	if err != nil {
		log.Fatal(err)
	}
	router := auth.SetupRouter(service)
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
