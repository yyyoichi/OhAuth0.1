package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
)

func main() {
	envPath := flag.String("source", "", "env file")
	flag.Parse()
	if envPath != nil {
		if err := godotenv.Load(*envPath); err != nil {
			panic(err)
		}
		log.Printf("read env file '%s'", *envPath)
	}

	var dbport string
	if dbport = os.Getenv("DATABASE_SERVER_PORT"); dbport == "" {
		panic("no required env found")
	}
	err := database.NewDatabaseServer(database.ServerConfig{
		Port: dbport,
	})
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	<-ctx.Done()
}
