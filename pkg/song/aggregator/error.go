package aggregator

import (
	"fmt"
)

type aggregationErr struct {
	msg     string
	reasons []error
}

func newAggregationErr(msg string, reasons ...error) *aggregationErr {
	return &aggregationErr{
		msg:     msg,
		reasons: reasons,
	}
}

func (err *aggregationErr) Error() string {
	return fmt.Sprintf("%s: %s", err.msg, err.reasons)
}
