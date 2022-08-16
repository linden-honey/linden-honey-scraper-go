package aggregator

import (
	"fmt"
)

type AggregationError struct {
	msg     string
	reasons []error
}

func NewAggregationError(msg string, reasons ...error) *AggregationError {
	return &AggregationError{
		msg:     msg,
		reasons: reasons,
	}
}

func (err *AggregationError) Error() string {
	return fmt.Sprintf("%s: %s", err.msg, err.reasons)
}
