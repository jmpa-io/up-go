package up

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// brokenReader simulates a body that always fails to read — used to test
// ErrSenderFailedParseResponse.
// from: https://stackoverflow.com/questions/45126312/how-do-i-test-an-error-on-reading-from-a-request-body
type brokenReader struct{}

func (br *brokenReader) Read(p []byte) (n int, err error) { return 0, fmt.Errorf("failed reading") }
func (br *brokenReader) Close() error                     { return fmt.Errorf("failed closing") }

// emptyErr is a placeholder error used to construct typed error strings for
// comparison in tests without caring about the inner message.
var emptyErr = fmt.Errorf("")

func Test_sender(t *testing.T) {
	tests := map[string]struct {
		mock    *mockRoundTripper
		request senderRequest
		result  interface{}
		err     string
	}{
		"catch json marshal error": {
			request: senderRequest{
				body: make(chan int), // channels cannot be marshalled to JSON
			},
			err: ErrFailedMarshal{emptyErr}.Error(),
		},
		"catch failed setup request": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{}
				},
			},
			request: senderRequest{
				method: "https://", // invalid method causes http.NewRequest to error
			},
			err: ErrSenderFailedSetupRequest{emptyErr}.Error(),
		},
		"catch failed send request": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{StatusCode: 401}
				},
			},
			request: senderRequest{
				method: "error", // triggers mockRoundTripper to return an error
			},
			err: ErrSenderFailedSendRequest{emptyErr}.Error(),
		},
		"catch failed parse response": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{Body: &brokenReader{}}
				},
			},
			err: ErrSenderFailedParseResponse{emptyErr}.Error(),
		},
		"catch json unmarshal error": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{}
				},
			},
			err: ErrFailedUnmarshal{emptyErr}.Error(),
		},
	}
	for name, tt := range tests {
		tt := tt
		ctx := context.Background()
		c := newTestClient(t, tt.mock)

		t.Run(name, func(t *testing.T) {
			_, err := c.sender(ctx, tt.request, &tt.result)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("sender() unexpected error;\nwant=%v\ngot=%v\n", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("sender() error;\n%v\n", err)
			}
		})
	}
}

