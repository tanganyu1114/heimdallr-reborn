package ginmiddleware

import "errors"

var (
	ErrClientNil     = errors.New("client cannot be nil")
	ErrApiGroupEmpty = errors.New("api group cannot be empty")
)
