package errtypes

import "errors"

var (
	NotFoundError      = errors.New("not found")
	MultipleFoundError = errors.New("multiple found")
)
