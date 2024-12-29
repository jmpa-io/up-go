package up

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

var (
	accountsTestdata1 = NewTestdata("accounts-1")
	accountsTestdata2 = NewTestdata("accounts-2")
	accountsTestdata3 = NewTestdata("accounts-3")
)

var location = time.FixedZone("AEST", 11*60*60)

func Test_GetAccounts(t *testing.T) {
	tests := map[string]struct {
		mock *MockRoundTripper
		want []AccountResource
		err  string
	}{
		"read accounts": {
			mock: &MockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					var b []byte
					switch {
					case strings.Contains(req.URL.String(), "---2"):
						b = accountsTestdata2.content
					case strings.Contains(req.URL.String(), "---3"):
						b = accountsTestdata3.content
					default:
						b = accountsTestdata1.content
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(b)),
						Header:     make(http.Header),
					}
				},
			},
			want: []AccountResource{
				{
					DisplayName:   "Spending",
					AccountType:   AccountTypeTransactional,
					OwnershipType: string(OwnershipTypeIndividual),
					Balance: Money{
						CurrencyCode:     "AUD",
						Value:            "1.00",
						ValueInBaseUnits: 100,
					},
					CreatedAt: time.Date(2024, 11, 06, 14, 26, 50, 00, location),
				},
				{
					DisplayName:   "Spending",
					AccountType:   AccountTypeTransactional,
					OwnershipType: string(OwnershipTypeIndividual),
					Balance: Money{
						CurrencyCode:     "AUD",
						Value:            "100.00",
						ValueInBaseUnits: 10000,
					},
					CreatedAt: time.Date(2024, 11, 06, 14, 26, 50, 00, location),
				},
				{
					DisplayName:   "Spending",
					AccountType:   AccountTypeTransactional,
					OwnershipType: string(OwnershipTypeIndividual),
					Balance: Money{
						CurrencyCode:     "AUD",
						Value:            "100.00",
						ValueInBaseUnits: 10000,
					},
					CreatedAt: time.Date(2024, 11, 06, 14, 26, 50, 00, location),
				},
			},
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
			got, err := c.ListAccounts(context.Background())
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf(
						"ListAccounts() returned an unexpected error;\nwant=%v\ngot=%v\n",
						tt.err,
						err,
					)
				}
				return
			}
			if err != nil {
				t.Errorf("ListAccounts() returned an error;\nerror=%v\n", err)
				return
			}

			// determine matches.
			matches := make(map[int]bool)
			for i := 0; i < len(got); i++ {
				g := got[i]
				w := tt.want[i] // TODO: could be an issue here with slice length.
				switch {
				case
					g.DisplayName == w.DisplayName,
					g.AccountType == w.AccountType,
					g.OwnershipType == w.OwnershipType,
					g.Balance == w.Balance,
					g.CreatedAt.Equal(w.CreatedAt):
					matches[i] = true
					continue
				}
			}
			var foundErrs bool
			for _, m := range matches {

				// skip passes!
				if m {
					continue
				}

				// fail on errors.
				foundErrs = true
			}
			if foundErrs {
				t.Errorf(
					"ListAccounts() returned unexpected configuration;\nwant=%+v\ngot=%+v\n",
					tt.want,
					got,
				)
			}
		})
	}
}
