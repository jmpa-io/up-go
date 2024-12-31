package up

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

var accountsTestdata []*testdata

func init() {

	// populate testdata.
	for i := 1; i <= 3; i++ {
		accountsTestdata = append(accountsTestdata, newTestdata(fmt.Sprintf("accounts-%v", i)))
	}
}

func Test_ListAccounts(t *testing.T) {
	tests := map[string]struct {
		mock *mockRoundTripper
		opts []ListAccountsOption
		want []AccountResource
		err  string
	}{
		"read accounts": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					b := accountsTestdata[0].content
					for i := 0; i < len(accountsTestdata); i++ {
						if !strings.Contains(req.URL.String(), fmt.Sprintf("---%v", i+1)) {
							continue
						}
						b = accountsTestdata[i].content
						break
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
					OwnershipType: AccountOwnershipTypeIndividual,
					Balance: Money{
						CurrencyCode:     "AUD",
						Value:            "1422.00",
						ValueInBaseUnits: 1422,
					},
					CreatedAt: time.Date(2024, 11, 06, 14, 26, 50, 00, location),
				},
				{
					DisplayName:   "Savings",
					AccountType:   AccountTypeSaver,
					OwnershipType: AccountOwnershipTypeIndividual,
					Balance: Money{
						CurrencyCode:     "AUD",
						Value:            "12455.00",
						ValueInBaseUnits: 12455,
					},
					CreatedAt: time.Date(2023, 04, 05, 06, 10, 20, 00, location),
				},
				{
					DisplayName:   "Home Loan",
					AccountType:   AccountTypeHomeLoan,
					OwnershipType: AccountOwnershipTypeJoint,
					Balance: Money{
						CurrencyCode:     "AUD",
						Value:            "52489.00",
						ValueInBaseUnits: 52489,
					},
					CreatedAt: time.Date(2022, 01, 14, 13, 30, 39, 00, location),
				},
			},
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
			got, err := c.ListAccounts(ctx, tt.opts...)

			// any errors?
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

			// do the lengths match?
			if len(got) != len(tt.want) {
				t.Errorf(
					"ListAccounts() returned unexpected number of results;\nwant=%d\ngot=%d\n",
					len(tt.want),
					len(got),
				)
				return
			}

			// is there a mismatch from what we're expecting vs what we've got?
			var foundErrs bool
			for i := 0; i < len(got); i++ {
				g := got[i]
				w := tt.want[i]
				if g.DisplayName != w.DisplayName ||
					g.AccountType != w.AccountType ||
					g.OwnershipType != w.OwnershipType ||
					g.Balance != w.Balance ||
					!g.CreatedAt.Equal(w.CreatedAt) {
					t.Errorf("mismatch at index %d;\nwant=%+v\ngot=%+v\n", i, w, g)
					foundErrs = true
				}
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
