package up

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func Test_sender(t *testing.T) {
	tests := map[string]struct {
		mock    *mockRoundTripper
		request senderRequest
		result  interface{}
		err     string
	}{
		"catch json marshal error": {
			request: senderRequest{
				data: make(chan int), // can't marshal a channel.
			},
			err: ErrFailedMarshal{fmt.Errorf("")}.Error(),
		},
		"catch failed setup request": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{}
				},
			},
			request: senderRequest{
				method: "https://", // returns an error on invalid method.
			},
			err: ErrSenderFailedSetupRequest{fmt.Errorf("")}.Error(),
		},
		"catch failed send request": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
					}
				},
			},
			request: senderRequest{
				method: "error", // not a real method, this is just to test.
			},
			err: ErrSenderFailedSendRequest{fmt.Errorf("")}.Error(),
		},
	}
	for name, tt := range tests {

		// tracing context.
		ctx := context.Background()

		// setup client with mock.
		c, _ := New(ctx, "xxxx",
			WithHttpClient(&http.Client{
				Transport: tt.mock,
			}),
		)

		// run tests.
		t.Run(name, func(t *testing.T) {
			resp, err := c.sender(ctx, tt.request, &tt.result)

			// any errors?
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf(
						"sender() returned an unexpected error;\nwant=%v\ngot=%v\n",
						tt.err,
						err,
					)
				}
				return
			}
			if err != nil {
				t.Errorf("sender() returned an error;\nerror=%v\n", err)
				return
			}

			fmt.Println(resp)
		})
	}
}
