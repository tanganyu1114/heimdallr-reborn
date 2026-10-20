package ginmiddleware

import "github.com/ClessLi/component-base/pkg/errors"

var (
	ErrClientNil     = errors.New("client cannot be nil")
	ErrApiGroupEmpty = errors.New("api group cannot be empty")
)
