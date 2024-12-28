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

	// send a ping to the API to check if the token is valid.
	p, err := c.Ping(ctx)
	if err != nil {
		fmt.Printf("failed to ping: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", p.Meta.StatusEmoji)

	// list accounts.
	accounts, err := c.ListAccounts(ctx)
	if err != nil {
		fmt.Printf("failed to list accounts: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", accounts)

	// list transactions.
	transactions, err := c.ListTransactions(ctx)
	if err != nil {
		fmt.Printf("failed to list transactions: %v\n", err)
	}
	fmt.Printf("%s\n", transactions)
}
