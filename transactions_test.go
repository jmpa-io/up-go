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

var transactionsTestdata []*testdata

func init() {

	// populate testdata.
	for i := 1; i <= 3; i++ {
		transactionsTestdata = append(
			transactionsTestdata,
			newTestdata(fmt.Sprintf("transactions-%v", i)),
		)
	}
}

func Test_ListTransactions(t *testing.T) {
	tests := map[string]struct {
		mock *mockRoundTripper
		want []TransactionResource
		err  string
	}{
		"read transactions": {
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					b := transactionsTestdata[0].content
					for i := 0; i < len(transactionsTestdata); i++ {
						if !strings.Contains(req.URL.String(), fmt.Sprintf("---%v", i+1)) {
							continue
						}
						b = transactionsTestdata[i].content
						break
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(b)),
						Header:     make(http.Header),
					}
				},
			},
			want: []TransactionResource{
				{
					Status:          "SETTLED",
					RawText:         "",
					Description:     "David Taylor",
					Message:         "Money for the pizzas last night.",
					IsCategorizable: true,
					Amount: Money{
						CurrencyCode:     "AUD",
						Value:            "-59.98",
						ValueInBaseUnits: -5998,
					},
					SettledAt: time.Date(2024, 11, 05, 07, 25, 12, 00, location),
					CreatedAt: time.Date(2024, 11, 05, 07, 25, 12, 00, location),
					PerformingCustomer: TransactionResourcePerformingCustomer{
						DisplayName: "Bobby",
					},
					DeepLinkURL: "up://transaction/VHJhbnNhY3Rpb24tMzg=",
				},
				{
					Status:          "HELD",
					RawText:         "SQ* TRAIN TICKETS",
					Description:     "John Doe",
					Message:         "Reimbursement for train tickets.",
					IsCategorizable: true,
					Amount: Money{
						CurrencyCode:     "AUD",
						Value:            "-20.50",
						ValueInBaseUnits: -2050,
					},
					CreatedAt: time.Date(2024, 11, 06, 8, 45, 00, 00, location),
					Note: TransactionResourceNote{
						Text: "Travel expense",
					},
					PerformingCustomer: TransactionResourcePerformingCustomer{
						DisplayName: "John",
					},
					DeepLinkURL: "up://transaction/VHJhbnNhY3Rpb24tNDk=",
				},
				{
					Status:          "SETTLED",
					RawText:         "",
					Description:     "Jane Smith",
					Message:         "Lunch meeting expense.",
					IsCategorizable: true,
					RoundUp: TransactionResourceRoundUp{
						Amount: Money{
							CurrencyCode:     "AUD",
							Value:            "-0.80",
							ValueInBaseUnits: -8,
						},
					},
					Amount: Money{
						CurrencyCode:     "AUD",
						Value:            "-45.30",
						ValueInBaseUnits: -4530,
					},
					SettledAt: time.Date(2024, 11, 7, 14, 15, 0, 0, location),
					CreatedAt: time.Date(2024, 11, 7, 14, 15, 0, 0, location),
					Note: TransactionResourceNote{
						Text: "Business lunch",
					},
					PerformingCustomer: TransactionResourcePerformingCustomer{
						DisplayName: "Jane",
					},
					DeepLinkURL: "up://transaction/VHJhbnNhY3Rpb24tNTY=",
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
			got, err := c.ListTransactions(ctx)

			// any errors?
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf(
						"ListTransactions() returned an unexpected error;\nwant=%v\ngot=%v\n",
						tt.err,
						err,
					)
				}
				return
			}
			if err != nil {
				t.Errorf("ListTransactions() returned an error;\nerror=%v\n", err)
				return
			}

			// do the lengths match?
			if len(got) != len(tt.want) {
				t.Errorf(
					"ListTransactions returned unexpected number of results;\nwant=%d\ngot=%d\n",
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
				if g.Status != w.Status ||
					g.RawText != w.RawText ||
					g.Description != w.Description ||
					g.Message != w.Message ||
					g.IsCategorizable != w.IsCategorizable ||
					g.HoldInfo != w.HoldInfo ||
					g.RoundUp != w.RoundUp ||
					g.Cashback != w.Cashback ||
					g.Amount != w.Amount ||
					g.ForeignAmount != w.ForeignAmount ||
					g.CardPurchaseMethod != w.CardPurchaseMethod ||
					!g.CreatedAt.Equal(w.CreatedAt) ||
					!g.SettledAt.Equal(w.SettledAt) ||
					g.TransactionType != w.TransactionType ||
					g.Note != w.Note ||
					g.PerformingCustomer != w.PerformingCustomer ||
					g.DeepLinkURL != w.DeepLinkURL {
					foundErrs = true
				}
			}
			if foundErrs {
				t.Errorf(
					"ListTransactions() returned unexpected configuration;\nwant=%+v\ngot=%+v\n",
					tt.want,
					got,
				)
			}
		})
	}
}
