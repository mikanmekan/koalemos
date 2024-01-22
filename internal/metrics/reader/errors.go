package reader

import (
	"errors"
)

var (
	ErrUnexpectedMetadata   = errors.New("unexpected metadata")
	ErrUnexpectedMetricLine = errors.New("unexpected metric line")
	ErrInvalidValue         = errors.New("invalid value in metric line")
	ErrOddLabelSetParts     = errors.New("odd number of label parts")
	ErrDuplicateLabelKey    = errors.New("label keys should not be repeated within a metric point")
)
