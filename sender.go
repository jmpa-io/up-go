package up

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// clientRequest simplifies sending a given request to the API by the sender.
type clientRequest struct {
	method  string
	path    string
	data    interface{}
	queries url.Values
}

// errorResponse represents an error returned from the API.
type errorResponse struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Source struct {
		Parameter string `json:"parameter"`
	}
}

// envelope defines errors returned from the API.
type envelope struct {
	Errors []errorResponse `json:"errors"`
}

// sender sends the given request to the API.s
func (c *Client) sender(cr clientRequest, result interface{}) (*http.Response, error) {

	// marshal body.
	var body []byte
	if !isNil(cr.data) {
		b, err := json.Marshal(cr.data)
		if err != nil {
			return nil, ErrFailedMarshal{err}
		}
		body = b
	}

	// setup request.
	req, err := http.NewRequest(cr.method, c.endpoint+cr.path, bytes.NewReader(body))
	if err != nil {
		return nil, ErrSenderFailedSetupRequest{err}
	}
	if cr.queries != nil {
		req.URL.RawQuery = cr.queries.Encode()
	}

	// add headers to request.
	req.Header = c.headers

	// send request.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, ErrSenderFailedSendRequest{err}
	}
	defer resp.Body.Close()

	// parse response.
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrSenderFailedParseResponse{err}
	}
	c.logger.Debug().
		Int("code", resp.StatusCode).
		Str("body", string(b)).
		Msg("response from API")

	// was this a valid request to the API?
	if http.StatusOK <= resp.StatusCode && resp.StatusCode < http.StatusMultipleChoices {
		if len(b) > 0 {
			return resp, json.Unmarshal(b, &result)
		}
		return resp, nil
	}

	// since we have an unexpected invalid response, return a generic response.
	var env envelope
	if err := json.Unmarshal(b, &env); err != nil {
		return nil, ErrFailedUnmarshal{err}
	}
	return nil, ErrSenderInvalidResponse{env, resp.StatusCode}
}
