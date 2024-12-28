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

func Test_GetTransactions(t *testing.T) {
	tests := map[string]struct {
		server   *httptest.Server
		testdata string
		want     *Ping
		err      string
	}{
		"successfully read transactions": {
			server: httptest.NewUnstartedServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					var out struct {
					}
					b, _ := json.Marshal(out)
					w.WriteHeader(http.StatusOK)
					w.Write(b)
				}),
			),
			testdata: "testdata/transactions.json",
		},
		// "unauthorized transaction": {
		// 	server: httptest.NewUnstartedServer(
		// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 			out := apiErrorResponse{
		// 				Errors: []apiErrorResponseError{
		// 					{
		// 						Status: "TODO",
		// 						Title:  "TODO",
		// 						Detail: "TODO",
		// 						Source: apiErrorResponseErrorSource{
		// 							Parameter: "TODO",
		// 						},
		// 					},
		// 				},
		// 			}
		// 			b, _ := json.Marshal(out)
		// 			w.WriteHeader(http.StatusUnauthorized)
		// 			w.Write(b)
		// 		}),
		// 	),
		// 	testdata: "testdata/unauthorized.json",
		// 	err:      "error response returned from API",
		// },
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
			got, err := c.ListTransactions(context.Background())
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf(
						"GetTransactions() returned an unexpected error; want=%v, got=%v",
						tt.err,
						err,
					)
				}
				return
			}
			if err != nil {
				t.Errorf("GetTransactions() returned an error; error=%v", err)
				return
			}
			switch {
			case
				!reflect.DeepEqual(got, tt.want):
				t.Errorf(
					"GetTransactions() returned unexpected configuration; want=%+v, got=%+v\n",
					tt.want,
					got,
				)
				return
			}
		})
	}
}
