package errors

import "errors"

var (
	ErrInternalServer = errors.New("internal error")
	ErrNotFound       = errors.New("resource not found")
)
