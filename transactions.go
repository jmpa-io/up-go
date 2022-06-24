package up

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HoldInfoObject struct {
	Amount        MoneyObject `json:"amount"`
	ForeignAmount MoneyObject `json:"foreignAmount"`
}

type RoundUpObject struct {
	Amount       MoneyObject `json:"amount"`
	BoostPortion MoneyObject `json:"boostPortion"`
}

type CashbackObject struct {
	Description string      `json:"description"`
	Amount      MoneyObject `json:"amount"`
}

type TransactionStatus string

const (
	TransactionStatusHeld    TransactionStatus = "HELD"
	TransactionStatusSettled TransactionStatus = "SETTELED"
)

type TransactionAttributesObject struct {
	Description   string            `json:"description"`
	Status        TransactionStatus `json:"status"`
	RawText       string            `json:"rawText"`
	Message       string            `json:"message"`
	HoldInfo      HoldInfoObject    `json:"holdInfo"`
	RoundUp       RoundUpObject     `json:"roundUp"`
	Cashback      CashbackObject    `json:"cashback"`
	Amount        MoneyObject       `json:"amount"`
	ForeignAmount MoneyObject       `json:"foreignAmount"`
	SettledAt     time.Time         `json:"settledAt"`
	CreatedAt     time.Time         `json:"createdAt"`
}

type TransactionRelationshipsObject struct {
	Account        AccountObject  `json:"account"`
	Tags           TagObject      `json:"tags"`
	ParentCategory CategoryObject `json:"parentCategory"`
	Category       CategoryObject `json:"category"`
}

type TransactionResourceObject struct {
	ResourceObject
	Attributes    TransactionAttributesObject    `json:"attributes"`
	Relationships TransactionRelationshipsObject `json:"relationships"`
}

// ListTransactionsResponse represents the response returned from the API when
// trying to list transactions for the authed account.
type ListTransactionsResponse struct {
	Data  []TransactionResourceObject `json:"data"`
	Links LinksObject                 `json:"links"`
}

// ListTransactionsOptions defines the options for listing transactions in the
// authed account.
type ListTransactionsOptions struct {
	name  string
	value string
}

func ListTransactionOptionPageSize(size int) ListTransactionsOptions {
	return ListTransactionsOptions{
		name:  "page[size]",
		value: strconv.Itoa(size),
	}
}

func ListTransactionOptionSince(since time.Time) ListTransactionsOptions {
	return ListTransactionsOptions{
		name:  "filter[since]",
		value: since.Format(time.RFC3339),
	}
}

func ListTransactionOptionUntil(until time.Time) ListTransactionsOptions {
	return ListTransactionsOptions{
		name:  "filter[until]",
		value: until.Format(time.RFC3339),
	}
}

// ListTransactions list all transactions for the authed account.
// https://developer.up.com.au/#get_transactions
func (c *Client) ListTransactions(options ...ListTransactionsOptions) ([]TransactionResourceObject, error) {

	// setup queries.
	queries := make(url.Values)
	for _, o := range options {
		queries[o.name] = []string{o.value}
	}

	// default queries.
	if _, ok := queries["page[size]"]; !ok {
		queries["page[size]"] = []string{"100"}
	}

	// setup request.
	cr := clientRequest{
		method:  http.MethodGet,
		path:    "/transactions",
		queries: queries,
	}

	// retrieve transactions.
	var transactions []TransactionResourceObject
	for {

		// get response.
		var resp ListTransactionsResponse
		if _, err := c.sender(cr, &resp); err != nil {
			return nil, err
		}

		// extract response data.
		transactions = append(transactions, resp.Data...)

		// paginate?
		if resp.Links.Next == "" {
			break
		}
		cr.path = strings.Replace(resp.Links.Next, c.endpoint, "", 1)
		cr.queries = nil
	}
	return transactions, nil
}
