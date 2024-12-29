package up

import "net/http"

// MockRoundTripper implements the RoundTripper interface.
type MockRoundTripper struct {
	MockFunc func(req *http.Request) *http.Response
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.MockFunc(req), nil
}
