package up

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func Test_setupQueries(t *testing.T) {
	tests := map[string]struct {
		options interface{}
		want    url.Values
	}{
		"setup queries for ListAccounts": {
			options: []ListAccountsOption{
				ListAccountsOptionPageSize(1),
				ListAccountsOptionFilterAccountOwnershipType(AccountOwnershipTypeIndividual),
				ListAccountsOptionFilterAccountType(AccountTypeTransactional),
			},
			want: url.Values{
				"page[size]":            []string{"1"},
				"filter[ownershipType]": []string{"INDIVIDUAL"},
				"filter[accountType]":   []string{"TRANSACTIONAL"},
			},
		},
		"setup queries for ListTags": {
			options: []ListTagsOption{
				ListTagsOptionPageSize(30),
			},
			want: url.Values{
				"page[size]": []string{"30"},
			},
		},
		"setup queries for ListTransactions": {
			options: []ListTransactionsOption{
				ListTransactionsOptionPageSize(10000),
				ListTransactionsOptionStatus(FilterStatusHeld),
				ListTransactionsOptionSince(time.Date(2022, 1, 10, 0, 0, 0, 0, location)),
				ListTransactionsOptionUntil(time.Date(2023, 12, 6, 0, 0, 0, 0, location)),
				ListTransactionsOptionCategory("hello"),
				ListTransactionsOptionTag("world"),
			},
			want: url.Values{
				"page[size]":       []string{"10000"},
				"filter[status]":   []string{"HELD"},
				"filter[since]":    []string{"2022-01-10T00:00:00+11:00"},
				"filter[until]":    []string{"2023-12-06T00:00:00+11:00"},
				"filter[category]": []string{"hello"},
				"filter[tag]":      []string{"world"},
			},
		},
		"check queries defaults page size": {
			options: []ListTagsOption{},
			want: url.Values{
				"page[size]": []string{"100"},
			},
		},
	}
	for name, tt := range tests {

		// run tests.
		t.Run(name, func(t *testing.T) {
			got := setupQueries(tt.options)

			// is there a mismatch from what we're expecting vs what we've got?
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"setupQueries() returned unexpected configuration;\nwant=%+v\ngot=%+v\n",
					tt.want,
					got,
				)
			}
		})
	}
}
