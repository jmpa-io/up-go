package up

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
)

// Ping represents a ping event returned from the API.
type Ping struct {
	Meta PingMeta `json:"meta"`
}

// PingMeta wraps the actual important values in the Ping response from the API.
type PingMeta struct {
	ID          string `json:"id"`
	StatusEmoji string `json:"statusEmoji"`
}

// Ping makes a ping request to the API. Use this to quickly determine if your
// provided token is working or not.
// https://developer.up.com.au/#get_util_ping.
func (c *Client) Ping(ctx context.Context) (*Ping, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "Ping")
	defer span.End()

	// ping.
	var p *Ping
	if _, err := c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   "/util/ping",
	}, &p); err != nil {
		return nil, err
	}
	return p, nil
}
