package up

import "net/http"

// Ping defines the response sent back from the Up API when sending a ping event.
type Ping struct {
	Id          string `json:"id"`
	StatusEmoji string `json:"statusEmoji"`
}

// Ping makes a ping request to the API.
// https://developer.up.com.au/#get_util_ping
func (c *Client) Ping() (*Ping, error) {
	cr := clientRequest{
		method: http.MethodGet,
		path:   "/util/ping",
	}
	var ping *Ping
	_, err := c.sender(cr, &ping)
	return ping, err
}
