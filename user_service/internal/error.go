package internal

const ErrorCodeValidation = "E400"

// NewValidationError returns a new Error with the given description and
// ErrorCodeValidation.
func NewValidationError(description string) error {
	return &Error{Code: ErrorCodeValidation, Description: description}
}

// Error implements error and provides a code and description.
type Error struct {
	Code, Description string
}

// Error implements the error interface.
func (e *Error) Error() string { return e.Description }
