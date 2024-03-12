package main

import (
	"context"
	"log"
	"os"

	"github.com/yyyoichi/OhAuth0.1/internal/database"
)

func main() {
	port := os.Getenv("DATABASE_SERVER_PORT")
	if port == "" {
		port = "3306"
	}
	err := database.NewDatabaseServer(database.ServerConfig{
		Port: port,
	})
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	<-ctx.Done()
}
