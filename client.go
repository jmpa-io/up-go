package up

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
)

// An iHttpClient is an interface over http.Client.
type iHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client defines a client for this package.
type Client struct {

	// tracing.
	tracerName string // The name of the tracer output in the traces.

	// config.
	endpoint      string      // The endpoint to query against.
	httpClient    iHttpClient // The http client used when sending / receiving data from the endpoint.
	headers       http.Header // The headers passed to the http client when sending / receiving data from the endpoint.
	skipAuthCheck bool        // Skip the Ping call on startup (useful when API is unreachable).

	// misc.
	logLevel slog.Level   // The log level of the default logger.
	logger   *slog.Logger // The logger used in this client (custom or default).
}

// New creates and returns a new Client, initialized with the provided token.
// The client itself is set up with tracing, logging, and HTTP configuration.
// Additional options can be provided to modify its behavior, via the options
// slice. The client is used for making requests and interacting with the Up
// Bank API.
func New(ctx context.Context, token string, options ...Option) (*Client, error) {

	// setup tracing.
	tracerName := "up-go"
	newCtx, span := otel.Tracer(tracerName).Start(ctx, "New")
	defer span.End()

	// check args.
	if token == "" {
		return nil, ErrClientEmptyToken{}
	}

	// clone http.DefaultTransport so we inherit all system-level settings
	// (DNS resolver, TLS config, proxy env vars) and only override the
	// connection-pool limits:
	//   • 30s timeout (DefaultClient has none — hangs forever)
	//   • MaxIdleConnsPerHost=20 so concurrent pagination workers can reuse
	//     keep-alive connections to the same host
	dt := http.DefaultTransport.(*http.Transport).Clone()
	dt.MaxIdleConns = 100
	dt.MaxIdleConnsPerHost = 20
	dt.IdleConnTimeout = 90 * time.Second
	c := &Client{
		tracerName: tracerName,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: dt,
		},
		endpoint: "https://api.up.com.au/api/v1",
	}

	// overwrite client with any given options.
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, ErrClientFailedToSetOption{err}
		}
	}

	// determine if the default logger should be used.
	if c.logger == nil {
		c.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: c.logLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05"))
				}
				return a
			},
		}))
	}

	// setup headers.
	headers := make(http.Header)
	headers.Set("Authorization", "Bearer "+token)
	headers.Set("Content-Type", "application/json")
	c.headers = headers

	// validate the token immediately by pinging the API.
	// Skipped if WithSkipAuthCheck was used.
	if !c.skipAuthCheck {
		if _, err := c.Ping(newCtx); err != nil {
			return nil, ErrClientFailedToPing{err}
		}
	}

	c.logger.Debug("client setup successfully")
	return c, nil
}

