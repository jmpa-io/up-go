package up

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	// test logger.
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05"))
			}
			return a
		},
	}))

	// pingMockTransport returns a successful ping response for any request —
	// used to prevent tests from hitting the real Up API during New() init.
	pingMockTransport = &mockRoundTripper{
		MockFunc: func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(pingTestdata.content)),
				Header:     make(http.Header),
			}
		},
	}

	// test httpClient with a mock transport so New() ping succeeds in tests.
	httpClient = &http.Client{
		Timeout:   30 * time.Second,
		Transport: pingMockTransport,
	}
)

func Test_New(t *testing.T) {
	tests := map[string]struct {
		token   string
		options []Option
		want    *Client
		err     string
	}{
		"default": {
			token: "xxxx",
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     slog.Default(),
			},
		},
		"no token": {
			want: &Client{},
			err:  "the provided token is empty",
		},
		"with log level (debug)": {
			token:   "xxxx",
			options: []Option{WithLogLevel(slog.LevelDebug)},
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     slog.Default(),
				logLevel:   -4,
			},
		},
		"with log level (info)": {
			token:   "xxxx",
			options: []Option{WithLogLevel(slog.LevelInfo)},
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     slog.Default(),
				logLevel:   0,
			},
		},
		"with log level (warn)": {
			token:   "xxxx",
			options: []Option{WithLogLevel(slog.LevelWarn)},
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     slog.Default(),
				logLevel:   4,
			},
		},
		"with log level (error)": {
			token:   "xxxx",
			options: []Option{WithLogLevel(slog.LevelError)},
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     slog.Default(),
				logLevel:   8,
			},
		},
		"with logger": {
			token:   "xxxx",
			options: []Option{WithLogger(logger)},
			want: &Client{
				httpClient: http.DefaultClient,
				logger:     logger,
			},
		},
		"with http client": {
			token:   "xxxx",
			options: []Option{WithHttpClient(httpClient)},
			want: &Client{
				httpClient: httpClient,
				logger:     slog.Default(),
			},
		},
	}
	for name, tt := range tests {

		// setup headers.
		headers := make(http.Header)
		if tt.token != "" {
			headers.Add("Authorization", "Bearer "+tt.token)
		}

		// inject a mock that returns a successful ping so New() doesn't hit
		// the real API during tests. We prepend it so test-specific options
		// (like WithHttpClient) can override it afterwards.
		pingMock := &mockRoundTripper{
			MockFunc: func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(pingTestdata.content)),
					Header:     make(http.Header),
				}
			},
		}
		// pingMock is prepended; test options override it if they supply their own httpClient.
		// For tests that supply WithHttpClient, we need the test's mock to also handle ping.
		// Solution: use pingMock unless the test already has WithHttpClient with a valid ping mock.
		opts := []Option{WithHttpClient(&http.Client{Transport: pingMock})}
		opts = append(opts, tt.options...)

		// run tests.
		t.Run(name, func(t *testing.T) {
			got, err := New(context.Background(), tt.token, opts...)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("New() returned an unexpected error; want=%v, got=%v", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("New() returned an error; error=%v", err)
				return
			}
			switch {
			case
				got.logLevel != tt.want.logLevel,
				(got.logger != slog.Default() && tt.want.logger != slog.Default()) && got.logger != tt.want.logger,
				got.headers.Get("Authorization") != headers.Get("Authorization"):
				t.Errorf(
					"New() returned unexpected configuration; want=%+v, got=%+v\n",
					tt.want,
					got,
				)
				return
			}
		})
	}
}
