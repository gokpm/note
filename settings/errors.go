package settings

import "errors"

var (
	ErrOption   = errors.New("invalid option")
	ErrFileName = errors.New("invalid file name")
	ErrEditor   = errors.New("invalid editor")
)
