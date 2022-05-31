package code

import "errors"

var (
	ErrNotFound = errors.New("ErrNotFound")
	NotFound    = 40001
	SystemErr   = 50001
)
