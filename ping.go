package up

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// Ping represents the response from the Up API ping endpoint.
type Ping struct {
	Meta PingMeta `json:"meta"`
}

// PingMeta contains the metadata returned by the ping endpoint.
type PingMeta struct {
	ID          string `json:"id"`
	StatusEmoji string `json:"statusEmoji"`
}

// Ping verifies that the token is valid and the Up API is reachable.
// It is called automatically during New() to validate the token.
// https://developer.up.com.au/#get_util_ping.
func (c *Client) Ping(ctx context.Context) (*Ping, error) {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "Ping")
	defer span.End()

	var p *Ping
	if _, err := c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   "/util/ping",
	}, &p); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("ping failed: %v", err))
		span.RecordError(err)
		return nil, fmt.Errorf("ping: %w", err)
	}
	return p, nil
}
