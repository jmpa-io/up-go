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

type AccountType string

const (
	AccountTypeSaver         AccountType = "SAVER"
	AccountTypeTransactional             = "TRANSACTIONAL"
	AccountTypeHomeLoan                  = "HOME_LOAN"
)

type OwnershipType string

const (
	OwnershipTypeIndividual OwnershipType = "INDIVIDUAL"
	OwnershipTypeJoint                    = "JOINT"
)

type AccountResource struct {
	DisplayName   string      `json:"displayName"`
	AccountType   AccountType `json:"accountType"`
	OwnershipType string      `json:"ownershipType"`
	Balance       Money       `json:"balance"`
	CreatedAt     time.Time   `json:"createdAt"`
}

type AccountRelationships struct {
	Transactions WrapperOnlyLinks `json:"transactions"`
}

// Account represents an account in Up.
type Account Data[AccountResource, AccountRelationships]

// AccountPaginationWrapper a pagination wrapper for a slice of AccountDataWrapper.
type AccountWrapper WrapperSlice[Account]

type ListAccountsOption struct {
	name  string
	value string
}

func ListAccountsOptionPageSize(size int) ListAccountsOption {
	return ListAccountsOption{
		name:  "page[size]",
		value: strconv.Itoa(size),
	}
}

func ListAccountsOptionFilterAccountType(t AccountType) ListAccountsOption {
	return ListAccountsOption{
		name:  "filter[accountType]",
		value: string(t),
	}
}

func ListAccountsOptionFilterOwnershipType(t OwnershipType) ListAccountsOption {
	return ListAccountsOption{
		name:  "filter[ownershipType]",
		value: string(t),
	}
}

// ListAccounts list all accounts for the authed account.
// https://developer.up.com.au/#get_accounts.
func (c *Client) ListAccounts(ctx context.Context,
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
		var resp AccountWrapper
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
