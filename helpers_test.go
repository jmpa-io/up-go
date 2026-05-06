package up

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
)

// newTestClient creates a Client for use in tests. The provided mock handles
// all requests AFTER the first one — the first request is always the init ping
// call, answered with a successful response so New() succeeds regardless of
// what the test mock does.
func newTestClient(t *testing.T, mock *mockRoundTripper) *Client {
	t.Helper()
	var callCount atomic.Int32
	combined := &mockRoundTripper{
		MockFunc: func(req *http.Request) *http.Response {
			if callCount.Add(1) == 1 {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(pingTestdata.content)),
					Header:     make(http.Header),
				}
			}
			if mock != nil {
				return mock.MockFunc(req)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(nil)),
				Header:     make(http.Header),
			}
		},
	}
	c, err := New(context.Background(), "xxxx", WithHttpClient(&http.Client{Transport: combined}))
	if err != nil {
		t.Fatalf("newTestClient: New() failed: %v", err)
	}
	return c
}

func Test_isNil(t *testing.T) {
	tests := map[string]struct {
		i    interface{}
		want bool
	}{
		"empty interface": {
			want: true,
		},
		"i has a string value": {
			i:    "hello world",
			want: false,
		},
		"i has a map value": {
			i: map[string]int{
				"hello": 1,
				"world": 2,
			},
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isNil(tt.i)
			if got != tt.want {
				t.Errorf("unexpected value returned; got: %v, want: %v\n", got, tt.want)
				return
			}
		})
	}
}

