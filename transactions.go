package up

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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
// to those who are either "HELD" or "SETTLED".
func ListTransactionsOptionStatus(status TransactionStatus) ListTransactionsOption {
	return ListTransactionsOption{newListOption("filter[status]", string(status))}
}

// ListTransactionsOptionSince filters the transactions returned from the API
// to those that occurred on or after the specified timestamp.
func ListTransactionsOptionSince(since time.Time) ListTransactionsOption {
	return ListTransactionsOption{newListOption("filter[since]", since.Format(time.RFC3339))}
}

// ListTransactionsOptionUntil filters the transactions returned from the API
// to those that occurred on or before the specified timestamp.
func ListTransactionsOptionUntil(until time.Time) ListTransactionsOption {
	return ListTransactionsOption{newListOption("filter[until]", until.Format(time.RFC3339))}
}

// ListTransactionsOptionCategory filters the transactions returned from the API
// to those that are categorized by the given Up category ID.
func ListTransactionsOptionCategory(category string) ListTransactionsOption {
	return ListTransactionsOption{newListOption("filter[category]", category)}
}

// ListTransactionsOptionTag filters the transactions returned from the API to
// those that are tagged with the given tag.
func ListTransactionsOptionTag(tag string) ListTransactionsOption {
	return ListTransactionsOption{newListOption("filter[tag]", tag)}
}

// listTransactions is the shared paginated fetch used by both ListTransactions
// and ListTransactionsByAccount.
func (c *Client) listTransactions(
	ctx context.Context,
	path string,
	options []ListTransactionsOption,
) (transactions []TransactionDataWrapper, err error) {

	sr := senderRequest{
		method:  http.MethodGet,
		path:    path,
		queries: setupQueries(options),
	}

	for {
		var resp TransactionsPaginationWrapper
		if _, err := c.sender(ctx, sr, &resp); err != nil {
			return nil, fmt.Errorf("fetching transactions from %s: %w", path, err)
		}
		transactions = append(transactions, resp.Data...)
		if resp.Links.Next == "" {
			break
		}
		sr.path = strings.Replace(resp.Links.Next, c.endpoint, "", 1)
		sr.queries = nil
	}
	return transactions, nil
}

// ListTransactions returns all transactions across all accounts for the
// authenticated user. Supports filtering via ListTransactionsOption.
// https://developer.up.com.au/#get_transactions.
func (c *Client) ListTransactions(ctx context.Context,
	options ...ListTransactionsOption,
) ([]TransactionDataWrapper, error) {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTransactions")
	defer span.End()

	txns, err := c.listTransactions(newCtx, "/transactions", options)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to list transactions: %v", err))
		span.RecordError(err)
		return nil, err
	}
	return txns, nil
}

// ListTransactionsByAccount returns all transactions for a specific account.
// https://developer.up.com.au/#get_accounts_accountId_transactions.
func (c *Client) ListTransactionsByAccount(
	ctx context.Context,
	accountID string,
	options ...ListTransactionsOption,
) ([]TransactionDataWrapper, error) {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListTransactionsByAccount")
	defer span.End()

	txns, err := c.listTransactions(newCtx,
		fmt.Sprintf("/accounts/%s/transactions", accountID),
		options,
	)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to list transactions for account %s: %v", accountID, err))
		span.RecordError(err)
		return nil, err
	}
	return txns, nil
}

// GetTransaction retrieves a single transaction by its ID, including the real
// timestamp, rawText, tags, and category relationship data.
// https://developer.up.com.au/#get_transactions_id.
func (c *Client) GetTransaction(
	ctx context.Context,
	id string,
) (*TransactionDataWrapper, error) {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetTransaction")
	defer span.End()

	var resp struct {
		Data TransactionDataWrapper `json:"data"`
	}
	if _, err := c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/transactions/%s", id),
	}, &resp); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get transaction %s: %v", id, err))
		span.RecordError(err)
		return nil, fmt.Errorf("getting transaction %s: %w", id, err)
	}
	return &resp.Data, nil
}

