// Package apperror defines the errors that services return to handlers.
//
// Each error has a Kind (one of the Err* values below) that decides the HTTP
// status, and a Message that is safe to show to the client:
//
//	return apperror.NotFound("manufacturer 7 not found")   // -> 404
//
// Callers check the kind with errors.Is:
//
//	if errors.Is(err, apperror.ErrNotFound) { ... }
package apperror

import "errors"

// Error kinds. httpx.WriteError maps each kind to an HTTP status.
var (
	ErrNotFound = errors.New("not found") // 404
	ErrConflict = errors.New("conflict")  // 409: duplicate, or still referenced by other rows
	// ErrInvalid is a well-formed request that breaks a validation or business rule (422).
	ErrInvalid = errors.New("invalid input")
	// ErrMalformed is a request that cannot be parsed: bad JSON, unknown field, wrong type (400).
	ErrMalformed = errors.New("malformed request")
	ErrTooLarge  = errors.New("request too large") // 413
)

// Error is an application error with a client-safe message.
type Error struct {
	Kind    error
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// Unwrap lets errors.Is(err, apperror.ErrNotFound) match on the Kind.
func (e *Error) Unwrap() error {
	return e.Kind
}

func NotFound(message string) error  { return &Error{Kind: ErrNotFound, Message: message} }
func Conflict(message string) error  { return &Error{Kind: ErrConflict, Message: message} }
func Invalid(message string) error   { return &Error{Kind: ErrInvalid, Message: message} }
func Malformed(message string) error { return &Error{Kind: ErrMalformed, Message: message} }
func TooLarge(message string) error  { return &Error{Kind: ErrTooLarge, Message: message} }
