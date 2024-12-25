package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmpa-io/up-go"
)

func main() {

	ctx := context.TODO()

	// retrieve token.
	token := os.Getenv("UP_TOKEN")

	// setup client.
	c, err := up.New(ctx, token, up.WithLogLevel(slog.LevelWarn))
	if err != nil {
		fmt.Printf("failed to setup client: %v\n", err)
		os.Exit(1)
	}

	// do something with the client..
	// like send a ping to the api to check if the token is valid.
	p, err := c.Ping(ctx)
	if err != nil {
		fmt.Printf("failed to ping: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", p.Meta.StatusEmoji)
}
