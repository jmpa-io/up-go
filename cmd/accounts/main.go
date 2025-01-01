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

	// list accounts.
	accounts, err := c.ListAccounts(ctx)
	if err != nil {
		fmt.Printf("failed to list accounts: %v\n", err)
		os.Exit(1)
	}
	for i, a := range accounts {
		fmt.Printf(
			"%v. %s, %s, %s, %v\n",
			i,
			a.DisplayName,
			a.OwnershipType,
			a.Balance.CurrencyCode,
			a.Balance.Value,
		)
	}
}
