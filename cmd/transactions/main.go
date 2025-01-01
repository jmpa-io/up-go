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

	// list transactions.
	transactions, err := c.ListTransactions(ctx)
	if err != nil {
		fmt.Printf("failed to list transactions: %v\n", err)
		os.Exit(1)
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
}
