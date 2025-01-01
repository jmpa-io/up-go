package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmpa-io/up-go"
)

func main() {

	// setup tracing.
	ctx := context.TODO()

	// retrieve token.
	token := os.Getenv("UP_TOKEN")

	// setup client.
	c, err := up.New(ctx, token, up.WithLogLevel(slog.LevelWarn))
	if err != nil {
		fmt.Printf("failed to setup client: %v\n", err)
		os.Exit(1)
	}

	// list tags.
	tags, err := c.ListTags(ctx)
	if err != nil {
		fmt.Printf("failed to list tags: %v\n", err)
		os.Exit(1)
	}
	for i, t := range tags {
		fmt.Printf("%v. %s %s\n", i, t.Type, t.ID)
	}
}
