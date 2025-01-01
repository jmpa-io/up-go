package up

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
)

// TransactionsPaginationWrapper is a pagination wrapper for a slice of
// TransactionDataWrapper. This type is used in organizing paginated
// transaction data received from the API.
type TransactionsPaginationWrapper WrapperSlice[TransactionDataWrapper]

// ListTransactionsOption defines the available options used to configure the
// ListTransactions function when listing transactions from the API.
type ListTransactionsOption struct {
	listOption
}

// ListTransactionsOptionPageSize sets the page size used when when listing
// transactions from the API. This option affects how many transactions are
// returned at once - increasing this can improve performance as the number of
// API calls is reduced.
func ListTransactionsOptionPageSize(size int) ListTransactionsOption {
	return ListTransactionsOption{newListOption("page[size]", strconv.Itoa(size))}
}

// ListTransactionsOptionStatus filters the transactions returned from the API
// to those who are either "HELD" or "SETTLED". Use this option if, for example,
// you'd like to list all transactions that are currently pending in Up.
func ListTransactionsOptionStatus(status TransactionStatus) ListTransactionsOption {
	return ListTransactionsOption{newListOption("filter[status]", string(status))}
}

// ListTransactionsOptionSince filters the transactions returned from the API
// to those that occurred on or after the specified 'since' timestamp. This is
// useful for retrieving transactions within a specific time window. When used
// with ListTransactionsOptionUntil, it defines the start of the window.
func ListTransactionsOptionSince(since time.Time) ListTransactionsOption {
	return ListTransactionsOption{newListOption("filter[since]", since.Format(time.RFC3339))}
}

// ListTransactionsOptionUntil filters the transactions returned from the API
// to those that occurred on or before the specified 'until' timestamp. When
// used with ListTransactionsOptionSince, it defines the end of the time
// window. If only ListTransactionsOptionUntil is provided, all transactions
// up to the specified time will be returned.
func ListTransactionsOptionUntil(until time.Time) ListTransactionsOption {
	return ListTransactionsOption{newListOption("filter[until]", until.Format(time.RFC3339))}
}

// ListTransactionsOptionCategory filters the transactions returned from the API
// to those that are categorized by the given category.
func ListTransactionsOptionCategory(category string) ListTransactionsOption {
	return ListTransactionsOption{newListOption("filter[category]", category)}
}

// ListTransactionsOptionTag filters the transactions returned from the API to
// those that are tagged with the given tag.
func ListTransactionsOptionTag(tag string) ListTransactionsOption {
	return ListTransactionsOption{newListOption("filter[tag]", tag)}
}

// ListTransactions returns a list of ALL transactions associated with authed
// user from the API. This function supports pagination, and is configurable by
// the given ListTransactionsOptions options.
// https://developer.up.com.au/#get_transactions.
func (c *Client) ListTransactions(ctx context.Context,
	options ...ListTransactionsOption,
) (transactions []TransactionResource, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTransactions")
	defer span.End()

	// setup request.
	sr := senderRequest{
		method:  http.MethodGet,
		path:    "/transactions",
		queries: setupQueries(options),
	}

	// retrieve transactions.
	for {

		// get response.
		var resp TransactionsPaginationWrapper
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
