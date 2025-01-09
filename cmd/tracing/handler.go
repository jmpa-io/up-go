package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"

	"github.com/jmpa-io/up-go"
)

type handler struct {

	// config.
	name        string
	version     string
	environment string

	// clients.
	upsvc *up.Client

	// misc.
	logger *slog.Logger
}

// run is like main but after the handler is configured.
func (h *handler) run(ctx context.Context) {

	// setup span.
	newCtx, span := otel.Tracer(h.name).Start(ctx, "run")
	defer span.End()

	// list accounts.
	accounts, err := h.upsvc.ListAccounts(newCtx)
	if err != nil {
		h.logger.Error("failed to list accounts", "error", err)
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
