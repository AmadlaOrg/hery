package entity

import "errors"

var (
	ErrorNotFound      = errors.New("not found")
	ErrorMultipleFound = errors.New("multiple found")
)
