package code

import "errors"

var (
	ErrNotFound = errors.New("ErrNotFound")
	Success     = 200
	NotFound    = 400
	SystemErr   = 500
)
