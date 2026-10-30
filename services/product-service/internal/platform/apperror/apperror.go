package apperror

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
	// ErrInvalid is a well-formed request that breaks a validation or business rule (422).
	ErrInvalid = errors.New("invalid input")
	// ErrMalformed is a request that cannot be parsed: bad JSON, unknown field, wrong type (400).
	ErrMalformed = errors.New("malformed request")
	ErrTooLarge  = errors.New("request too large")
)

// Error carries a client-safe message while still matching a sentinel via errors.Is.
type Error struct {
	Kind    error
	Message string
}

func (e *Error) Error() string { return e.Message }
func (e *Error) Unwrap() error { return e.Kind }

func NotFound(msg string) error { return &Error{Kind: ErrNotFound, Message: msg} }
func Conflict(msg string) error { return &Error{Kind: ErrConflict, Message: msg} }
func Invalid(msg string) error  { return &Error{Kind: ErrInvalid, Message: msg} }
func Malformed(msg string) error {
	return &Error{Kind: ErrMalformed, Message: msg}
}
func TooLarge(msg string) error { return &Error{Kind: ErrTooLarge, Message: msg} }
