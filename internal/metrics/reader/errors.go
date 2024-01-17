package reader

import (
	"errors"
)

var (
	ErrUnexpectedMetadata   = errors.New("unexpected metadata")
	ErrUnexpectedMetricLine = errors.New("unexpected metric line")
	ErrOddLabelSetParts     = errors.New("odd number of label parts")
	ErrDuplicateLabelKey    = errors.New("label keys should not be repeated")
)
