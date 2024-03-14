package main

import (
	"context"
	"log"

	"github.com/yyyoichi/OhAuth0.1/internal/resource"
)

func main() {
	ctx := context.Background()
	service, err := resource.NewService(ctx, resource.Config{
		DatabaseServerURL: "http://localhost:3306",
	})
	if err != nil {
		log.Fatal(err)
	}
	router := resource.SetupRouter(service)
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
