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

var accountTestdata = newTestdata("account")

func Test_GetAccount(t *testing.T) {
	tests := map[string]struct {
		mock *mockRoundTripper
		id   string
		want *AccountResource
		err  string
	}{
		"get account": {
			id: "9bcc759f-0ff8-458e-abd5-59bd47c9b560",
			mock: &mockRoundTripper{
				MockFunc: func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBuffer(accountTestdata.content)),
						Header:     make(http.Header),
					}
				},
			},
			want: &AccountResource{
				DisplayName:   "🐷 Savings",
				AccountType:   AccountTypeSaver,
				OwnershipType: AccountOwnershipTypeIndividual,
				Balance: Money{
					CurrencyCode:     "AUD",
					Value:            "125.36",
					ValueInBaseUnits: 12536,
				},
				CreatedAt: time.Date(2024, 11, 06, 14, 26, 51, 00, location),
			},
		},
	}
	for name, tt := range tests {
		ctx := context.Background()
		c := newTestClient(t, tt.mock)
		t.Run(name, func(t *testing.T) {
			got, err := c.GetAccount(ctx, tt.id)
			if tt.err != "" && err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("GetAccount() returned an unexpected error;\nwant=%v\ngot=%v\n", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Errorf("GetAccount() returned an error;\nerror=%v\n", err)
				return
			}
			if got.DisplayName != tt.want.DisplayName ||
				got.AccountType != tt.want.AccountType ||
				got.OwnershipType != tt.want.OwnershipType ||
				got.Balance != tt.want.Balance ||
				!got.CreatedAt.Equal(tt.want.CreatedAt) {
				t.Errorf("GetAccount() returned unexpected result;\nwant=%+v\ngot=%+v\n", tt.want, got)
			}
		})
	}
}
