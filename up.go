package up

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// based on https://github.com/nick96/upngo/blob/master/upngo.go

const (
	// endpoint is the default endpoint for Up Banks's API.
	endpoint = "https://api.up.com.au/api/v1"
	// defaultTimeout is the default timeout duration used on HTTP requests.
	defaultTimeout = time.Second * 300
	// defaultCode is the default error code for failures.
	defaultCode = -1
)

// Error defines an error received when making a request to the API.
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error returns a string representing the error, satisfying the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("upbank: %s (%d)", e.Message, e.Code)
}

// UpBank defines the UpBank client.
type UpBank struct {
	httpClient *http.Client // the HTTP client used to query against the API.
	token      string       // the UpBank API token.
}

// NewClient returns a new UpBank API client which can be used to make
// requests to the UpBank API.
func NewClient(token string) *UpBank {
	return &UpBank{
		token:      token,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
}

// clientRequest defines information that can be used to make a request to UpBank.
type clientRequest struct {
	Method  string      `json:"-"`
	Path    string      `json:"-"`
	Queries url.Values  `json:"-"`
	Data    interface{} `json:"data,omitempty"`
}

type errorResponse struct {
	Errors []struct {
		Status string `json:"status"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Source struct {
			Parameter string `json:"parameter"`
		}
	} `json:"errors"`
}

func (e errorResponse) Error() string {
	var out string
	for _, err := range e.Errors {
		if out != "" {
			out += "; "
		}
		out += fmt.Sprintf("%s %s: %s", err.Status, err.Title, err.Detail)
	}
	return out
}

// request makes a request to UpBank's API.
func (ub *UpBank) request(cr clientRequest, result interface{}) (*http.Response, error) {

	var err error

	// Marshal body.
	var body []byte
	if cr.Data != nil {
		body, err = json.Marshal(cr.Data)
		if err != nil {
			return nil, Error{fmt.Sprintf("Failed to marshal JSON: %s", err), defaultCode}
		}
	}

	// Construct request.
	cr.Path = strings.Replace(cr.Path, endpoint, "", -1)
	req, err := http.NewRequest(cr.Method, endpoint+cr.Path, bytes.NewReader(body))
	if err != nil {
		return nil, Error{fmt.Sprintf("Failed to construct request: %s", err), defaultCode}
	}
	req.Header.Set("Authorization", "Bearer "+ub.token)
	if cr.Queries != nil {
		req.URL.RawQuery = cr.Queries.Encode()
	}

	// fmt.Println(req.URL)

	// Make request.
	resp, err := ub.httpClient.Do(req)
	if err != nil {
		return nil, Error{fmt.Sprintf("Failed to make request: %s", err), defaultCode}
	}
	defer resp.Body.Close()

	// Parse response.
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, Error{fmt.Sprintf("Failed to read response: %d %s", resp.StatusCode, err), defaultCode}
	}

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		if len(b) == 0 {
			return resp, nil
		}
		return resp, json.Unmarshal(b, &result)
	}

	// parse error.
	var errResp errorResponse
	if err := json.Unmarshal(b, &errResp); err != nil {
		return nil, Error{fmt.Sprintf("Failed to parse error response: %d %s ", resp.StatusCode, err), defaultCode}
	}
	return nil, errResp
}

// TransactionsOption is an option for the transactions API.
type TransactionsOption struct {
	name  string
	value string
}

// TODO: This is a horrible name but otherwise it clashes with the accounts page
// size option. Anyway, I'm thinking either have to test out whether we can have
// a base option then alias it for the specific ones but I'm not sure if that
// will still be typesafe. Alternatively, it might be good, in general, to split
// out each resoure into a subclient, then they could have different namespaces.
func WithTransactionPageSize(size int) TransactionsOption {
	return TransactionsOption{
		name:  "page[size]",
		value: strconv.Itoa(size),
	}
}

func WithFilterSince(since time.Time) TransactionsOption {
	return TransactionsOption{
		name:  "filter[since]",
		value: since.Format(time.RFC3339),
	}
}

func WithFilterUntil(until time.Time) TransactionsOption {
	return TransactionsOption{
		name:  "filter[until]",
		value: until.Format(time.RFC3339),
	}
}

// ListTransactions list all transactions in the account.
// https://developer.up.com.au/#get_transactions
func (ub *UpBank) ListTransactions(options ...TransactionsOption) ([]TransactionsResponse, error) {
	queries := make(url.Values)
	for _, option := range options {
		queries[option.name] = []string{option.value}
	}
	cr := clientRequest{
		Method:  http.MethodGet,
		Path:    "/transactions",
		Queries: queries,
	}
	var transactions []TransactionsResponse
	for {
		tr := TransactionsResponse{}
		if _, err := ub.request(cr, &tr); err != nil {
			return nil, err
		}
		transactions = append(transactions, tr)
		if tr.Links.Next == "" {
			break
		}
		cr.Path = tr.Links.Next
		cr.Queries = nil
	}
	return transactions, nil
}

type PingResponse struct {
	Id          string `json:"id"`
	StatusEmoji string `json:"statusEmoji"`
}

// Ping to quickly check auth.
// https://developer.up.com.au/#get_util_ping
func (ub *UpBank) Ping() (PingResponse, error) {
	cr := clientRequest{
		Method: http.MethodGet,
		Path:   "/util/ping",
	}
	pr := PingResponse{}
	_, err := ub.request(cr, &pr)
	return pr, err
}

type TagLabels struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// AddTagsToTransactions ...
// https://developer.up.com.au/#post_transactions_transactionId_relationships_tags
func (ub *UpBank) AddTagsToTransactions(id string, labels []TagLabels) error {
	cr := clientRequest{
		Method: http.MethodPost,
		Path:   fmt.Sprintf("/transactions/%s/relationships/tags", id),
		Data:   labels,
	}
	var r interface{}
	_, err := ub.request(cr, &r)
	return err
}
