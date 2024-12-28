package up

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel"
)

// senderRequest simplifies sending a given request to the API endpoint, by the sender.
type senderRequest struct {
	method  string
	path    string
	data    interface{}
	queries url.Values
}

type apiErrorResponseErrorSource struct {
	Parameter string `json:"parameter"`
}

type apiErrorResponseError struct {
	Status string                      `json:"status"`
	Title  string                      `json:"title"`
	Detail string                      `json:"detail"`
	Source apiErrorResponseErrorSource `json:"source"`
}

// apiErrorResponse represents one or more errors returned from the API.
type apiErrorResponse struct {
	Errors []apiErrorResponseError `json:"errors"`
}

// sender sends the given senderRequest to the API endpoint.
func (c *Client) sender(
	ctx context.Context,
	sr senderRequest,
	result interface{},
) (resp *http.Response, err error) {

	// setup tracing.
	_, span := otel.Tracer(c.tracerName).Start(ctx, "sender")
	defer span.End()

	// marshal body.
	var body []byte
	if !isNil(sr.data) {
		body, err = json.Marshal(sr.data)
		if err != nil {
			return nil, ErrFailedMarshal{err}
		}
	}

	// setup request.
	req, err := http.NewRequest(sr.method, c.endpoint+sr.path, bytes.NewReader(body))
	if err != nil {
		return nil, ErrSenderFailedSetupRequest{err}
	}
	if sr.queries != nil {
		req.URL.RawQuery = sr.queries.Encode()
	}

	// add headers to request.
	req.Header = c.headers

	// send request.
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, ErrSenderFailedSendRequest{err}
	}
	defer resp.Body.Close()

	// parse response.
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrSenderFailedParseResponse{err}
	}

	// determine if the response was successful or a failure.
	if http.StatusOK <= resp.StatusCode && resp.StatusCode < http.StatusMultipleChoices {
		c.logger.Debug("response from API", "code", resp.StatusCode, "body", string(b))
		if len(b) > 0 {
			return resp, json.Unmarshal(b, &result)
		}
		return resp, nil
	}

	c.logger.Error("response from API", "code", resp.StatusCode, "body", string(b))
	var errs apiErrorResponse
	if err := json.Unmarshal(b, &errs); err != nil {
		return nil, ErrFailedUnmarshal{err}
	}
	return nil, ErrSenderInvalidResponse{errs, resp.StatusCode}
}

// ---

// ErrSenderFailedSetupRequest is returned whenever the sender fails to
// setup a new http request.
type ErrSenderFailedSetupRequest struct {
	err error
}

func (e ErrSenderFailedSetupRequest) Error() string {
	return fmt.Sprintf("failed to setup http request: %v", e.err)
}

// ErrSenderFailedSendRequest is returned whenever the sender fails to send
// a new http request.
type ErrSenderFailedSendRequest struct {
	err error
}

func (e ErrSenderFailedSendRequest) Error() string {
	return fmt.Sprintf("failed to send http request: %v", e.err)
}

// ErrSenderFailedParseResponse is returned when the sender fails to parse a
// response from the API.
type ErrSenderFailedParseResponse struct {
	err error
}

func (e ErrSenderFailedParseResponse) Error() string {
	return fmt.Sprintf("failed to parse response: %v", e.err)
}

// ErrSenderInvalidResponse is returned when the sender receives an error
// response specifically from the API.
type ErrSenderInvalidResponse struct {
	errs       apiErrorResponse
	statusCode int
}

func (e ErrSenderInvalidResponse) Error() string {
	var errs []string
	for _, err := range e.errs.Errors {
		errs = append(errs, err.Title)
	}
	return fmt.Sprintf(
		"error response returned from API; status_code=%v, count=%v, errors=%s",
		e.statusCode,
		len(errs),
		strings.Join(errs, ";"),
	)
}
