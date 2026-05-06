package up

import "fmt"

// ErrClientEmptyToken is returned when no token is provided to the client.
type ErrClientEmptyToken struct {
}

func (e ErrClientEmptyToken) Error() string {
	return "the provided token is empty"
}

// ErrClientFailedToPing is returned when the ping request to validate the
// token fails during client initialisation.
type ErrClientFailedToPing struct {
	err error
}

func (e ErrClientFailedToPing) Error() string {
	return fmt.Sprintf("failed to ping Up API (token may be invalid): %v", e.err)
}

// ErrClientFailedToSetOption is returned when an option encounters an error
// when trying to be set with the client.
type ErrClientFailedToSetOption struct {
	err error
}

func (e ErrClientFailedToSetOption) Error() string {
	return fmt.Sprintf("failed to set option in client: %v", e.err)
}
