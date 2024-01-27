package publisher

import (
	"context"
	"io"
)

type Service interface {
	Publish(ctx context.Context, in io.Reader) error
}
