package up

import "log/slog"

// Option configures a departure client.
type Option func(*Client) error

// WithLogger overwrites the default logger in this client with a given custom logger.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

// WithLogLevel configures the log level for the default logger in this client.
func WithLogLevel(level slog.Level) Option {
	return func(c *Client) error {
		c.logLevel = level
		return nil
	}
}

// WithHttpClient overirdes the default httpClient used when sending /
// receieving data from the endpoint.
func WithHttpClient(httpClient iHttpClient) Option {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}
