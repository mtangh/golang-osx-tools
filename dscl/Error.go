/* Error.go */

package dscl

// Error ...
type Error struct {
	exitStatus int
	stdErr     []byte
}

func newDsclError(exitStatus int, stdErr []byte) *Error {
	return &Error{exitStatus, stdErr}
}

func (dsclErr *Error) Error() string {
	return string(dsclErr.stdErr)
}
