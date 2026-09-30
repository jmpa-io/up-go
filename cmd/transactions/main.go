package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmpa-io/up-go"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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
		a := t.Attributes
		fmt.Printf(
			"%v. %s, %s, %s, %s, %s, %v, %s, %v\n",
			i,
			a.CreatedAt,
			a.SettledAt,
			a.Status,
			a.Amount.CurrencyCode,
			a.CardPurchaseMethod.Method,
			a.Amount.Value,
			a.RawText,
			a.RoundUp.Amount.Value,
		)
	}
}
