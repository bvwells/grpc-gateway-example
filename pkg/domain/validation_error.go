package domain

// NewValidationError returns a new ValidationError.
func NewValidationError(msg string) ValidationError {
	return ValidationError{msg: msg}
}

// ValidationError is a validation error.
type ValidationError struct {
	msg string
}

// Error returns the validation error string.
func (e ValidationError) Error() string {
	return e.msg
}
