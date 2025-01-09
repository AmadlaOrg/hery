package file

import "errors"

var (
	ErrorFailedToStatFile = errors.New("failed to stat file")
	ErrorNotAFile         = errors.New("not a file")
	ErrorMagicIsNotAMatch = errors.New("magic is not a match")
)
