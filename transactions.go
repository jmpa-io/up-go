package up

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
)

type TransactionAttributes struct {
	Status             string      `json:"status"`
	RawText            string      `json:"rawText"`
	Description        string      `json:"description"`
	Message            string      `json:"message"`
	IsCategorizable    bool        `json:"isCategorizable"`
	HoldInfo           interface{} `json:"holdInfo"`
	RoundUp            interface{} `json:"roundUp"`
	Cashback           interface{} `json:"cashback"`
	Amount             Amount      `json:"amount"`
	ForeignAmount      interface{} `json:"foreignAmount"`
	CardPurchaseMethod interface{} `json:"cardPurchaseMethod"`
	SettledAt          time.Time   `json:"settledAt"`
	CreatedAt          time.Time   `json:"createdAt"`
	TransactionType    interface{} `json:"transactionType"`
	Note               string      `json:"note"`
	PerformingCustomer struct {
		DisplayName string `json:"displayName"`
	} `json:"performingCustomer"`
	DeepLinkURL string `json:"deepLinkURL"`
}

type TransactionRelationships struct {
	Account struct {
		Data struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
		Links RelatedLink `json:"links"`
	} `json:"account"`
	TransferAccount interface{} `json:"transferAccount"`
	Category        interface{} `json:"category"`
	ParentCategory  interface{} `json:"parentCategory"`
	Tags            TagData     `json:"tags"`
	Attachment      interface{} `json:"attachment"`
}

// TransactionData represents a transaction in Up.
type TransactionDataWrapper Data[TransactionAttributes, TransactionRelationships]

// TransactionsWrapper is a pagination wrapper for a slice of TransactionData.
type TransactionPaginationWrapper PaginationWrapper[TransactionDataWrapper]

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
func (c *Client) ListTransactions(ctx context.Context,
	options ...ListTransactionsOptions,
) ([]TransactionAttributes, error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTransactions")
	defer span.End()

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
	sr := senderRequest{
		method:  http.MethodGet,
		path:    "/transactions",
		queries: queries,
	}

	// retrieve transactions response.
	var transactions []TransactionAttributes
	for {

		// get response.
		var resp TransactionPaginationWrapper
		if _, err := c.sender(newCtx, sr, &resp); err != nil {
			return nil, err
		}

		// extract response data.
		for _, t := range resp.Data {
			transactions = append(transactions, t.Attributes)
		}

		// paginate?
		if resp.Links.Next == "" {
			break
		}
		sr.path = strings.Replace(resp.Links.Next, c.endpoint, "", 1)
		sr.queries = nil
	}
	return transactions, nil
}
