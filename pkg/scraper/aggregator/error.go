package aggregator

import (
	"fmt"
)

// AggregationError represents an aggregation error object.
type AggregationError struct {
	msg     string
	reasons []error
}

// NewAggregationError return a pointer to a new instance of the aggregation error.
func NewAggregationError(msg string, reasons ...error) *AggregationError {
	return &AggregationError{
		msg:     msg,
		reasons: reasons,
	}
}

// Error returns an aggregated error message.
func (err *AggregationError) Error() string {
	return fmt.Sprintf("%s: %s", err.msg, err.reasons)
}
