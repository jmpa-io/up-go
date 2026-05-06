package up

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// AccountsPaginationWrapper is a pagination wrapper for a slice of AccountDataWrapper.
type AccountsPaginationWrapper WrapperSlice[AccountDataWrapper]

// ListAccountsOption configures a ListAccounts call.
type ListAccountsOption struct {
	listOption
}

// ListAccountsOptionPageSize sets the number of accounts returned per page.
func ListAccountsOptionPageSize(size int) ListAccountsOption {
	return ListAccountsOption{newListOption("page[size]", strconv.Itoa(size))}
}

// ListAccountsOptionFilterAccountType filters accounts to SAVER, TRANSACTIONAL,
// or HOME_LOAN.
func ListAccountsOptionFilterAccountType(t AccountType) ListAccountsOption {
	return ListAccountsOption{newListOption("filter[accountType]", string(t))}
}

// ListAccountsOptionFilterAccountOwnershipType filters accounts to INDIVIDUAL
// or JOINT ownership.
func ListAccountsOptionFilterAccountOwnershipType(t AccountOwnershipType) ListAccountsOption {
	return ListAccountsOption{newListOption("filter[ownershipType]", string(t))}
}

// ListAccounts returns all accounts for the authenticated user.
// https://developer.up.com.au/#get_accounts.
func (c *Client) ListAccounts(
	ctx context.Context,
	opts ...ListAccountsOption,
) (accounts []AccountResource, err error) {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAccounts")
	defer span.End()

	sr := senderRequest{
		method:  http.MethodGet,
		path:    "/accounts",
		queries: setupQueries(opts),
	}

	for {
		var resp AccountsPaginationWrapper
		if _, err := c.sender(newCtx, sr, &resp); err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("failed to list accounts: %v", err))
			span.RecordError(err)
			return nil, fmt.Errorf("listing accounts: %w", err)
		}
		for _, a := range resp.Data {
			accounts = append(accounts, a.Attributes)
		}
		if resp.Links.Next == "" {
			break
		}
		sr.path = strings.Replace(resp.Links.Next, c.endpoint, "", 1)
		sr.queries = nil
	}

	return accounts, nil
}

// GetAccount retrieves a single account by its ID.
// https://developer.up.com.au/#get_accounts_id.
func (c *Client) GetAccount(ctx context.Context, id string) (*AccountResource, error) {

	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "GetAccount")
	defer span.End()

	var resp struct {
		Data AccountDataWrapper `json:"data"`
	}
	if _, err := c.sender(newCtx, senderRequest{
		method: http.MethodGet,
		path:   fmt.Sprintf("/accounts/%s", id),
	}, &resp); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf("failed to get account %s: %v", id, err))
		span.RecordError(err)
		return nil, fmt.Errorf("getting account %s: %w", id, err)
	}
	attrs := resp.Data.Attributes
	return &attrs, nil
}
