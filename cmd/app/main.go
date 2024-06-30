package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	serviceclient "github.com/yyyoichi/OhAuth0.1/internal/service-client"
)

func main() {
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

	var srcport string
	if srcport = os.Getenv("RESOURCE_SERVER_PORT"); srcport == "" {
		panic("no required env found")
	}

	var authport string
	if authport = os.Getenv("AUTHORIZATION_SERVER_PORT"); authport == "" {
		panic("no required env found")
	}
	var uiport string
	if uiport = os.Getenv("UI_SERVER_PORT"); uiport == "" {
		panic("no required env found")
	}

	sc := bufio.NewScanner(os.Stdin)
	brawser := serviceclient.NewBrawser(serviceclient.BrawserConfig{
		RedirectPort:      7777,
		AuthServerURI:     "http://localhost:" + authport,
		ResourceServerURI: "http://localhost:" + srcport,
		AuthUIURI:         "http://localhost:" + uiport + "/v1/auth",
	})
	ctx := context.Background()
	go func() {
		for {
			fmt.Printf("\nPlease enter the command... \n")
			sc.Scan()
			input := sc.Text()
			output, err := brawser.Brawse(ctx, input)
			if err != nil {
				fmt.Printf("\nError!!: %s", err.Error())
				continue
			}
			fmt.Printf("\n%s", output.Msg())
		}
	}()
	<-ctx.Done()
}
