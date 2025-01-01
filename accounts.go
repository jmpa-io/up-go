package up

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
)

// AccountsPaginationWrapper is a pagination wrapperfor a slice of AccountData.
// This type is used in organizing paginated account data received from the API.
type AccountsPaginationWrapper WrapperSlice[AccountDataWrapper]

// ListAccountsOption defines the available options used to configure the
// ListAccounts function when listing accounts from the API.
type ListAccountsOption struct {
	listOption
}

// ListAccountsOptionPageSize sets the page size used when when listing
// accounts from the API. This option affects how many accounts are returned at
// once - increasing this can improve performance as the number of API calls
// is reduced.
func ListAccountsOptionPageSize(size int) ListAccountsOption {
	return ListAccountsOption{newListOption("page[size]", strconv.Itoa(size))}
}

// ListAccountsOptionFilterAccountType filters the accounts returned from the
// API to those who are either "SAVER", "TRANSACTIONAL", or "HOME_LOAN". Use
// this option if, for example, you'd like to list all accounts that are
// transactional and not associated with your savings or your home loan.
func ListAccountsOptionFilterAccountType(t AccountType) ListAccountsOption {
	return ListAccountsOption{newListOption("filter[accountType]", string(t))}
}

// ListAccountsOptionFilterAccountOwnershipType filters the accounts returned
// from the API to those who are either "INDIVIDUAL" or "JOINT". Use this option
// if, for example, you'd like to list all accounts that belong to you only..
func ListAccountsOptionFilterAccountOwnershipType(t AccountOwnershipType) ListAccountsOption {
	return ListAccountsOption{newListOption("filter[ownershipType]", string(t))}
}

// ListAccounts returns a list of ALL accounts associated with authed user from
// the API. This function supports pagination, and is configurable by the given
// ListAccountsOptions options.
// https://developer.up.com.au/#get_accounts.
func (c *Client) ListAccounts(
	ctx context.Context,
	opts ...ListAccountsOption,
) (accounts []AccountResource, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAccounts")
	defer span.End()

	// setup request.
	sr := senderRequest{
		method:  http.MethodGet,
		path:    "/accounts",
		queries: setupQueries(opts),
	}

	// retrieve accounts.
	for {

		// get response.
		var resp AccountsPaginationWrapper
		if _, err := c.sender(newCtx, sr, &resp); err != nil {
			return nil, err
		}

		// extract response data.
		for _, a := range resp.Data {
			accounts = append(accounts, a.Attributes)
		}

		// paginate?
		if resp.Links.Next == "" {
			break
		}
		sr.path = strings.Replace(resp.Links.Next, c.endpoint, "", 1)
		sr.queries = nil
	}

	return accounts, nil
}
