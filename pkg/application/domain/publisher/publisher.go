package publisher

import (
	"context"
	"io"
)

// Service is an interface of publisher use-cases.
type Service interface {
	Publish(ctx context.Context, artifactName string, in io.Reader) error
}
