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

	// list transactions.
	transactions, err := c.ListTransactions(ctx)
	if err != nil {
		fmt.Printf("failed to list transactions: %v\n", err)
	}
	for i, t := range transactions {
		fmt.Printf(
			"%v. %s, %s, %s, %s, %s, %v, %s, %v\n",
			i,
			t.CreatedAt,
			t.SettledAt,
			t.Status,
			t.Amount.CurrencyCode,
			t.CardPurchaseMethod.Method,
			t.Amount.Value,
			t.RawText,
			t.RoundUp.Amount.Value,
		)
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
