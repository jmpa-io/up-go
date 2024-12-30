package up

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
)

// AccountsPaginationWrapper a pagination wrapper for a slice of AccountDataWrapper.
type AccountsPaginationWrapper WrapperSlice[Account]

type ListAccountsOption struct {
	ListOption
}

func ListAccountsOptionPageSize(size int) ListAccountsOption {
	return ListAccountsOption{NewListOption("page[size]", strconv.Itoa(size))}
}

func ListAccountsOptionFilterAccountType(t AccountType) ListAccountsOption {
	return ListAccountsOption{NewListOption("filter[accountType]", string(t))}
}

func ListAccountsOptionFilterAccountOwnershipType(t AccountOwnershipType) ListAccountsOption {
	return ListAccountsOption{NewListOption("filter[ownershipType]", string(t))}
}

// ListAccounts list ALL accounts for the authed user.
// https://developer.up.com.au/#get_accounts.
func (c *Client) ListAccounts(
	ctx context.Context,
	options ...ListAccountsOption,
) (accounts []AccountResource, err error) {

	// setup tracing.
	newCtx, span := otel.Tracer(c.tracerName).Start(ctx, "ListAccounts")
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
		path:    "/accounts",
		queries: queries,
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
