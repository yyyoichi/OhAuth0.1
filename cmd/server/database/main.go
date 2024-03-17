package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/yyyoichi/OhAuth0.1/internal/database"
)

func main() {
	ctx := context.Background()
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(l)

	envPath := flag.String("source", "", "env file")
	flag.Parse()
	if envPath != nil {
		if err := godotenv.Load(*envPath); err != nil {
			panic(err)
		}
		slog.Info("read env file", slog.String("path", *envPath))
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
	<-ctx.Done()
}
