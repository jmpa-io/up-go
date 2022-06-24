package up

import "net/http"

// Ping represents a ping event sent from the API.
type Ping struct {
	Id          string `json:"id"`
	StatusEmoji string `json:"statusEmoji"`
}

// PingResponse represents the response returned from the API when
// trying to ping the API.
type PingResponse struct {
	Meta Ping `json:"meta"`
}

// Ping makes a ping request to the API.
// https://developer.up.com.au/#get_util_ping
func (c *Client) Ping() (*Ping, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   "/util/ping",
	}
	var resp *PingResponse
	_, err := c.sender(cr, &resp)
	return &resp.Meta, err
}
