package reader

import (
	"errors"
)

var (
	ErrUnexpectedMetadata   = errors.New("unexpected metadata")
	ErrUnexpectedMetricLine = errors.New("unexpected metric line")
)
