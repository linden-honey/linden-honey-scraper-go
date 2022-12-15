package aggregator

import (
	"fmt"
)

// AggregationError is an error with multiple causes and a general message.
type AggregationError struct {
	msg    string
	causes []error
}

// NewAggregationError return a pointer to the new instance of [AggregationError].
func NewAggregationError(msg string, causes ...error) *AggregationError {
	return &AggregationError{
		msg:    msg,
		causes: causes,
	}
}

// Error returns an aggregated error message.
func (err *AggregationError) Error() string {
	return fmt.Sprintf("%s: %s", err.msg, err.causes)
}
