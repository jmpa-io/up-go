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

// transactionsTestdata is declared in transactions_test.go.

var transactionTestdata = newTestdata("transaction")

func Test_GetTransaction(t *testing.T) {
	tests := map[string]struct {
		mock *mockRoundTripper
		id   string
		want *TransactionDataWrapper
		err  string
	}{
		"get transaction": {
			id: "6fff09f5-be7d-4ae1-9f71-4a25440bc405",
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(transactionTestdata.content)),
						Header:     make(http.Header),
					}
				},
			},
			want: &TransactionDataWrapper{
				Attributes: TransactionResource{
					Status:          TransactionStatusSettled,
					RawText:         "WARUNG BEBEK, UBUD INDONES",
					Description:     "Warung Bebek Bengil",
					IsCategorizable: true,
					Amount: Money{
						CurrencyCode:     "AUD",
						Value:            "-107.92",
						ValueInBaseUnits: -10792,
					},
					CreatedAt: time.Date(2024, 11, 03, 04, 00, 00, 00, location),
					SettledAt: time.Date(2024, 11, 03, 04, 00, 00, 00, location),
				},
			},
		},
	}
	for name, tt := range tests {
		ctx := context.Background()
		c := newTestClient(t, tt.mock)
		t.Run(name, func(t *testing.T) {
			got, err := c.GetTransaction(ctx, tt.id)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("GetTransaction() returned an unexpected error;\nwant=%v\ngot=%v\n", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("GetTransaction() returned an error;\nerror=%v\n", err)
				return
			}
			w := tt.want.Attributes
			g := got.Attributes
			if g.Status != w.Status ||
				g.RawText != w.RawText ||
				g.Description != w.Description ||
				g.IsCategorizable != w.IsCategorizable ||
				g.Amount != w.Amount ||
				!g.CreatedAt.Equal(w.CreatedAt) ||
				!g.SettledAt.Equal(w.SettledAt) {
				t.Errorf("GetTransaction() returned unexpected result;\nwant=%+v\ngot=%+v\n", w, g)
			}
		})
	}
}

func Test_ListTransactionsByAccount(t *testing.T) {
	tests := map[string]struct {
		mock      *mockRoundTripper
		accountID string
		wantCount int
		err       string
	}{
		"list transactions by account": {
			accountID: "fe848390-7d39-41fb-b01d-545de29ab74b",
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
			wantCount: 3,
		},
	}
	for name, tt := range tests {
		ctx := context.Background()
		c := newTestClient(t, tt.mock)
		t.Run(name, func(t *testing.T) {
			got, err := c.ListTransactionsByAccount(ctx, tt.accountID)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("ListTransactionsByAccount() returned an unexpected error;\nwant=%v\ngot=%v\n", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("ListTransactionsByAccount() returned an error;\nerror=%v\n", err)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("ListTransactionsByAccount() returned unexpected count;\nwant=%d\ngot=%d\n", tt.wantCount, len(got))
			}
		})
	}
}
