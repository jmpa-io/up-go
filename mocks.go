package up

import "net/http"

// mockRoundTripper implements the RoundTripper interface.
type mockRoundTripper struct {
	MockFunc func(req *http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.MockFunc(req), nil
}
