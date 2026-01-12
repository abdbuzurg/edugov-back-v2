// Package apperr defines structured application errors with HTTP-friendly codes.
package apperr

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Code identifies the application error category.
type Code string

const (
	CodeValidation   Code = "VALIDATION"
	CodeNotFound     Code = "NOT_FOUND"
	CodeConflict     Code = "CONFLICT"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeUnavailable  Code = "UNAVAILABLE"
	CodeInternal     Code = "INTERNAL"
)

// Error is a domain/app error:
// - msg is safe to return to client
// - cause is internal (logged only)
type Error struct {
	code   Code
	msg    string
	fields map[string]string
	cause  error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.cause == nil {
		return fmt.Sprintf("%s: %s", e.code, e.msg)
	}
	return fmt.Sprintf("%s: %s: %v", e.code, e.msg, e.cause)
}

// Unwrap exposes the underlying cause, if any.
func (e *Error) Unwrap() error { return e.cause }

// Code returns the error classification.
func (e *Error) Code() Code { return e.code }
// Msg returns the client-safe message.
func (e *Error) Msg() string { return e.msg }
// Fields returns a copy of field-level validation errors.
func (e *Error) Fields() map[string]string { return cloneMap(e.fields) }

func cloneMap(in map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// As extracts an *Error from a wrapped error chain.
func As(err error) (*Error, bool) {
	var ae *Error
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

// Validation builds a validation error with optional field-level details.
func Validation(msg string, fields map[string]string) error {
	return &Error{code: CodeValidation, msg: msg, fields: cloneMap(fields)}
}
// NotFound builds a not-found error.
func NotFound(msg string) error { return &Error{code: CodeNotFound, msg: msg} }
// Conflict builds a conflict error.
func Conflict(msg string) error { return &Error{code: CodeConflict, msg: msg} }
// Unauthorized builds an unauthorized error.
func Unauthorized(msg string) error { return &Error{code: CodeUnauthorized, msg: msg} }
// Forbidden builds a forbidden error.
func Forbidden(msg string) error { return &Error{code: CodeForbidden, msg: msg} }

// Unavailable builds a service-unavailable error.
func Unavailable(msg string, cause error) error {
	if msg == "" {
		msg = "service unavailable"
	}
	return &Error{code: CodeUnavailable, msg: msg, cause: cause}
}

// Internal builds an internal error with an optional cause.
func Internal(msg string, cause error) error {
	if msg == "" {
		msg = "internal error"
	}
	return &Error{code: CodeInternal, msg: msg, cause: cause}
}

// ValidationFromValidator converts validator errors into Validation errors.
func ValidationFromValidator(err error) error {
	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make(map[string]string, len(ve))
		for _, fe := range ve {
			// fe.Field() is struct field name; fe.Tag() is validation tag (e.g. "required", "email")
			// You can customize messages per tag:
			fields[fe.Field()] = humanizeValidation(fe)
		}
		return Validation("validation failed", fields)
	}

	// not a validation error => treat as internal
	return Internal("internal error", fmt.Errorf("validator: %w", err))
}

func humanizeValidation(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required"
	case "email":
		return "invalid email"
	case "min":
		return "too short"
	case "max":
		return "too long"
	default:
		// fallback: include tag
		return fe.Tag()
	}
}
