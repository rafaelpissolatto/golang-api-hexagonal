package custom_error

// StatusError is a status custom implementation of error.
type StatusError struct {
	statusCode int
	message    string
}

// New returns an error with status code and message
func New(statusCode int, message string) error {
	return &StatusError{statusCode, message}
}

// Error default func that return the error message
func (e *StatusError) Error() string {
	return e.message
}

// ErrorCode return the error status code
func (e *StatusError) ErrorCode() int {
	return e.statusCode
}
