package up

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
)

// TransactionsWrapper is a pagination wrapper for a slice of TransactionData.
type TransactionWrapper WrapperSlice[Transaction]

// ListTransactionsOption defines the options for listing transactions in the
// authed account.
type ListTransactionsOption struct {
	ListOption
}

type FilterStatus string

const (
	FilterStatusHeld    FilterStatus = "HELD"
	FilterStatusSettled FilterStatus = "SETTLED"
)

func ListTransactionsOptionPageSize(size int) ListTransactionsOption {
	return ListTransactionsOption{NewListOption("page[size]", strconv.Itoa(size))}
}

func ListTransactionsOptionStatus(status FilterStatus) ListTransactionsOption {
	return ListTransactionsOption{NewListOption("filter[status]", string(status))}
}

func ListTransactionsOptionSince(since time.Time) ListTransactionsOption {
	return ListTransactionsOption{NewListOption("filter[since]", since.Format(time.RFC3339))}
}

func ListTransactionsOptionUntil(until time.Time) ListTransactionsOption {
	return ListTransactionsOption{NewListOption("filter[until]", until.Format(time.RFC3339))}
}

func ListTransactionsOptionCategory(category string) ListTransactionsOption {
	return ListTransactionsOption{NewListOption("filter[category]", category)}
}

func ListTransactionsOptionTag(tag string) ListTransactionsOption {
	return ListTransactionsOption{NewListOption("filter[tag]", tag)}
}

// ListTransactions list all transactions for the authed account.
// https://developer.up.com.au/#get_transactions
func (c *Client) ListTransactions(ctx context.Context,
	opts ...ListTransactionsOption,
) (transactions []TransactionResource, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTransactions")
	defer span.End()

	// setup request.
	sr := senderRequest{
		method:  http.MethodGet,
		path:    "/transactions",
		queries: setupQueries(opts),
	}

	// retrieve transactions.
	for {

		// get response.
		var resp TransactionWrapper
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
		c.logger.Debug(resp.Links.Next)
		sr.path = strings.Replace(resp.Links.Next, c.endpoint, "", 1)
		sr.queries = nil
	}
	return transactions, nil
}
