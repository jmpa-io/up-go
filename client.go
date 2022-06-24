package up

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

// An iHttpClient is an interface over http.Client.
type iHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client defines an Up Bank client; the interface for this package.
type Client struct {
	logLevel   LogLevel
	httpClient iHttpClient
	headers    http.Header
	endpoint   string // the endpoint to query against.

	// misc.
	logger zerolog.Logger
}

// New returns a client for this package, which can be used to make
// requests to the Up Bank API.
func New(token string, options ...Option) (*Client, error) {

	// check args.
	if token == "" {
		return nil, ErrMissingToken{}
	}

	// default client.
	c := &Client{
		logLevel:   LogLevelDebug,
		httpClient: http.DefaultClient,
		endpoint:   "https://api.up.com.au/api/v1",
	}

	// overwrite client with any given options.
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, ErrFailedOptionSet{err}
		}
	}

	// setup logger.
	zerolog.MessageFieldName = "msg"
	var level zerolog.Level
	switch c.logLevel {
	case LogLevelDebug:
		level = zerolog.DebugLevel
	case LogLevelInfo:
		level = zerolog.InfoLevel
	case LogLevelWarn:
		level = zerolog.WarnLevel
	case LogLevelError:
		level = zerolog.ErrorLevel
	default:
		level = zerolog.ErrorLevel
	}
	c.logger = zerolog.New(os.Stderr).
		With().Caller().Logger().
		With().Timestamp().Logger().Level(level)
	c.logger.Debug().Msg("setting up client")

	// setup headers.
	headers := make(http.Header)
	headers.Set("Authorization", "Bearer "+token)
	headers.Set("Content-Type", "application/json")
	c.headers = headers

	c.logger.Debug().Msg("client setup successfully")
	return c, nil
}
