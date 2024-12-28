package up

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func Test_Ping(t *testing.T) {
	tests := map[string]struct {
		server   *httptest.Server
		testdata string
		want     *Ping
		err      string
	}{
		"successful ping": {
			server: httptest.NewUnstartedServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					out := Ping{
						Meta: PingMeta{
							ID:          "14101d6e-03ee-4623-849c-c320e0764649",
							StatusEmoji: "⚡️",
						},
					}
					b, _ := json.Marshal(out)
					w.WriteHeader(http.StatusOK)
					w.Write(b)
				}),
			),
			testdata: "testdata/ping.json",
		},
		"unauthorized ping": {
			server: httptest.NewUnstartedServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					out := apiErrorResponse{
						Errors: []apiErrorResponseError{
							{
								Status: "TODO",
								Title:  "TODO",
								Detail: "TODO",
								Source: apiErrorResponseErrorSource{
									Parameter: "TODO",
								},
							},
						},
					}
					b, _ := json.Marshal(out)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(b)
				}),
			),
			testdata: "testdata/unauthorized.json",
			err:      "error response returned from API",
		},
	}
	for name, tt := range tests {

		// tracing context.
		ctx := context.Background()

		// start test server.
		tt.server.Start()
		defer tt.server.Close()

		// read + parse testdata.
		if err := readTestdata(tt.testdata, &tt.want); err != nil {
			panic(err)
		}

		// setup client with test server.
		c, _ := New(ctx, "xxxx",
			WithHttpClient(tt.server.Client()),
			WithEndpoint(tt.server.URL),
		)

		// run tests.
		t.Run(name, func(t *testing.T) {
			got, err := c.Ping(context.Background())
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
				!reflect.DeepEqual(got, tt.want):
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
