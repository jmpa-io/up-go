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
	pingTestdata         = newTestdata("ping")
	unauthorizedTestdata = newTestdata("unauthorized")
)

func Test_Ping(t *testing.T) {
	tests := map[string]struct {
		// responseCode and responseBody control what the mock returns for
		// the actual Ping call inside the test — not for New() init.
		responseCode int
		responseBody []byte
		want         *Ping
		err          string
	}{
		"successful ping": {
			responseCode: http.StatusOK,
			responseBody: pingTestdata.content,
			want: &Ping{
				Meta: PingMeta{
					ID:          "14101d6e-03ee-4623-849c-c320e0764649",
					StatusEmoji: "⚡️",
				},
			},
		},
		"unauthorized ping": {
			responseCode: http.StatusUnauthorized,
			responseBody: unauthorizedTestdata.content,
			err:          "error response returned from API",
		},
	}
	for name, tt := range tests {
		tt := tt
		ctx := context.Background()

		// Create a counter so the mock returns success for the init ping
		// (first call) and the test-specific response for subsequent calls.
		callCount := 0
		mock := &mockRoundTripper{
			MockFunc: func(req *http.Request) *http.Response {
				callCount++
				if callCount == 1 {
					// first call is the init ping — always succeed
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(pingTestdata.content)),
						Header:     make(http.Header),
					}
				}
				return &http.Response{
					StatusCode: tt.responseCode,
					Body:       io.NopCloser(bytes.NewBuffer(tt.responseBody)),
					Header:     make(http.Header),
				}
			},
		}

		c, err := New(ctx, "xxxx", WithHttpClient(&http.Client{Transport: mock}))
		if err != nil {
			t.Errorf("Test_Ping/%s: New() unexpectedly failed: %v", name, err)
			continue
		}

		t.Run(name, func(t *testing.T) {
			got, err := c.Ping(ctx)
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ping() returned unexpected configuration;\nwant=%+v\ngot=%+v\n", tt.want, got)
			}
		})
	}
}
