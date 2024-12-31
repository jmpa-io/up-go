package up

import (
	"fmt"
	"net/http"
	"strings"
)

// mockRoundTripper implements the RoundTripper interface.
type mockRoundTripper struct {
	MockFunc func(req *http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.Contains(req.Method, "error") {
		return nil, fmt.Errorf("an error occurred")
	}
	return m.MockFunc(req), nil
}
