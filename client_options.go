package up

import (
	"log/slog"
	"net/http"
)

// Option configures a departure client.
type Option func(*Client) error

// WithLogLevel sets the log level for the default logger.
func WithLogLevel(level slog.Level) Option {
	return func(c *Client) error {
		c.logLevel = level
		return nil
	}
}

// WithLogger overwrites the default logger with the given custom logger.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

// WithHttpClient overwrites the default httpClient used for API communication.
// Pass an *http.Client directly, or wrap a custom http.RoundTripper using
// &http.Client{Transport: rt}.
func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithSkipAuthCheck disables the Ping call that normally happens at startup.
// Use this when the Up API may be unreachable at init time (e.g. restrictive
// corporate networks). Tools will still fail gracefully if unreachable.
func WithSkipAuthCheck() Option {
	return func(c *Client) error {
		c.skipAuthCheck = true
		return nil
	}
}

// WithEndpoint overrides the default Up API base URL. Useful for testing
// against a mock server or a future staging environment.
func WithEndpoint(endpoint string) Option {
	return func(c *Client) error {
		c.endpoint = endpoint
		return nil
	}
}

// WithTracerName overrides the tracer name used in OpenTelemetry spans.
func WithTracerName(name string) Option {
	return func(c *Client) error {
		c.tracerName = name
		return nil
	}
}
