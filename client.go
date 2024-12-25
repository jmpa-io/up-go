package up

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
)

// An iHttpClient is an interface over http.Client.
type iHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client defines an Up Bank client; the interface for this package.
type Client struct {

	// tracing.
	tracerName string // The name of the tracer output in the traces.

	// config.
	endpoint   string      // The endpoint to query against.
	httpClient iHttpClient // The http client used when sending / receiving data from the endpoint.
	headers    http.Header // The headers passed to the http client when sending / receiving data from the endpoint.

	// misc.
	logLevel slog.Level   // The log level of the default logger.
	logger   *slog.Logger // The logger used in this client (custom or default).
}

// New returns a client for this package, which can be used to make
// requests to the Up Bank API.
func New(ctx context.Context, token string, options ...Option) (*Client, error) {

	// setup tracing.
	tracerName := "up-go"
	_, span := otel.Tracer(tracerName).Start(ctx, "New")
	defer span.End()

	// check args.
	if token == "" {
		return nil, ErrClientEmptyToken{}
	}

	// default client.
	c := &Client{
		httpClient: http.DefaultClient,
		endpoint:   "https://api.up.com.au/api/v1",
	}

	// overwrite client with any given options.
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, ErrClientFailedToSetOption{err}
		}
	}

	// determine if the default logger should be used.
	if c.logger == nil {

		// use default logger.
		c.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: c.logLevel, // default log level is 'INFO'.
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

	c.logger.Debug("client setup successfully")
	return c, nil
}

// ---

// ErrClientEmptyToken is returned when no token is provided to the client.
type ErrClientEmptyToken struct {
}

func (e ErrClientEmptyToken) Error() string {
	return "the provided token is empty"
}

// ErrClientFailedToSetOption is returned when an option encounters an error
// when trying to be set with the client.
type ErrClientFailedToSetOption struct {
	err error
}

func (e ErrClientFailedToSetOption) Error() string {
	return fmt.Sprintf("failed to set option in client: %v", e.err)
}
