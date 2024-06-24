package apprerrors

import (
	"errors"
	"fmt"
	"net/http"
)

type Type string

type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

const (
	Authorization      = "AUTHORIZATION"
	ServiceUnavailable = "SERVICEUNAVAILABLE"

	BadRequest      = "BADREQUEST"
	Conflict        = "CONFLICT"
	Internal        = "INTERNAL"
	NotFound        = "NOTFOUND"
	PayloadTooLarge = "PAYLOADTOOLARGE"
)

func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		{
			return http.StatusUnauthorized
		}
	case BadRequest:
		{
			return http.StatusBadRequest
		}
	case Conflict:
		{
			return http.StatusConflict
		}
	case ServiceUnavailable:
		{
			return http.StatusServiceUnavailable
		}
	case Internal:
		{
			return http.StatusInternalServerError
		}
	case NotFound:
		{
			return http.StatusNotFound
		}
	case PayloadTooLarge:
		{
			return http.StatusRequestEntityTooLarge
		}
	default:
		return http.StatusInternalServerError
	}

}

func Status(e error) int {
	var customErr *Error
	if errors.As(e, &customErr) {
		return customErr.Status()
	}
	return http.StatusInternalServerError
}

func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}
func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Bad request. REason :: %v  ", reason),
	}
}
func NewConflict(name, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("Conflict for :: %v  with :: %v", name, value),
	}
}
func NewTimedOut() *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: fmt.Sprintf("Request is timed out, seems liek service is unavailable"),
	}
}
func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: "Internal error.",
	}
}
func NewNotFound(name, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("Resource for :: %v not found with ::: %v", name, value),
	}
}
func NewPayloadTooLarge(actualSize, maxSize int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Body size limit in :: %d exceeded. Actual payload size ::: %d", maxSize, actualSize),
	}
}
