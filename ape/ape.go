package ape

import (
	"errors"
)

type Error struct {
	// id unique error identifier
	// in uppercase format like "ADMIN_CAN_NOT_DELETE_SELF"
	id string

	// internal error which caused this error
	cause error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return e.id
}

func (e *Error) Is(target error) bool {
	var be *Error
	if errors.As(target, &be) {
		return e.id == be.id
	}
	return false
}

func (e *Error) Raise(cause error) error {
	return &Error{
		id:    e.id,
		cause: cause,
	}
}

func DeclareError(id string) *Error {
	return &Error{
		id: id,
	}
}
