package publisher

import (
	"context"
	"io"
)

// Service is an interface of publisher use-cases.
type Service interface {
	Publish(ctx context.Context, name string, in io.Reader) error
}
