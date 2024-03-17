package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/yyyoichi/OhAuth0.1/internal/auth"
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

	service, err := auth.NewService(ctx, auth.Config{
		DatabaseServerURL: "http://localhost:" + dbport,
	})
	if err != nil {
		log.Fatal(err)
	}

	var port string
	if port = os.Getenv("AUTHORIZATION_SERVER_PORT"); port == "" {
		panic("no required env found")
	}
	var uiport string
	if uiport = os.Getenv("UI_SERVER_PORT"); uiport == "" {
		panic("no required env found")
	}
	router := auth.SetupRouter(service, fmt.Sprintf("http://localhost:%s", uiport))
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
