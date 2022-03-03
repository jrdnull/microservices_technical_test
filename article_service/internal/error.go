package internal

import "fmt"

// Internal error codes.
const (
	ErrorCodeValidation   = "E400"
	ErrorCodeUnauthorized = "E401"
)

// ErrUnauthorized is an Error with Code set to ErrorCodeUnauthorized.
var ErrUnauthorized = &Error{Code: ErrorCodeUnauthorized, Description: "unauthorized"}

// NewValidationError returns a new Error with the given description and
// ErrorCodeValidation.
func NewValidationError(description string) error {
	return &Error{Code: ErrorCodeValidation, Description: description}
}

// NewValidationErrorf returns a new Error with the given formatted description
// and ErrorCodeValidationf.
func NewValidationErrorf(format string, a ...interface{}) error {
	return &Error{Code: ErrorCodeValidation, Description: fmt.Sprintf(format, a...)}
}

// Error implements error and provides a code and description.
type Error struct {
	Code, Description string
}

// Error implements the error interface.
func (e *Error) Error() string { return e.Description }
