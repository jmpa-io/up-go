package up

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var (
	pingTestdata         = NewTestdata("ping")
	unauthorizedTestdata = NewTestdata("unauthorized")
)

func Test_Ping(t *testing.T) {
	tests := map[string]struct {
		mock *MockRoundTripper
		want *Ping
		err  string
	}{
		"successful ping": {
			mock: &MockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(pingTestdata.content)),
					}
				},
			},
			want: &Ping{
				Meta: PingMeta{
					ID:          "14101d6e-03ee-4623-849c-c320e0764649",
					StatusEmoji: "⚡️",
				},
			},
		},
		"unauthorized ping": {
			mock: &MockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       io.NopCloser(bytes.NewBuffer(unauthorizedTestdata.content)),
					}
				},
			},
			err: "error response returned from API",
		},
	}
	for name, tt := range tests {

		// tracing context.
		ctx := context.Background()

		// setup client with test server.
		c, _ := New(ctx, "xxxx",
			WithHttpClient(&http.Client{
				Transport: tt.mock,
			}),
		)

		// run tests.
		t.Run(name, func(t *testing.T) {
			got, err := c.Ping(context.Background())
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("Ping() returned an unexpected error;\nwant=%v\ngot=%v\n", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Ping() returned an error;\nerror=%v\n", err)
				return
			}
			switch {
			case
				!reflect.DeepEqual(got, tt.want):
				t.Errorf(
					"Ping() returned unexpected configuration;\nwant=%+v\ngot=%+v\n",
					tt.want,
					got,
				)
				return
			}
		})
	}
}
