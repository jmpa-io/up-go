package up

import "fmt"

// ErrFailedMarshal is returned when an error is returned from json.Marshal.
type ErrFailedMarshal struct {
	err error
}

func (e ErrFailedMarshal) Error() string {
	return fmt.Sprintf("failed to marshal data: %v", e.err)
}

// ErrFailedUnmarshal is returned when an error returned from json.Unmarshal.
type ErrFailedUnmarshal struct {
	err error
}

func (e ErrFailedUnmarshal) Error() string {
	return fmt.Sprintf("failed to unmarshal data: %v", e.err)
}
