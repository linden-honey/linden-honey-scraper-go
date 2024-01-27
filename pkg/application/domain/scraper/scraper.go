package scraper

import (
	"context"
	"io"
)

// Service is an interface of scraper use-cases.
type Service interface {
	ScrapeSongs(ctx context.Context, out io.Writer) error
}
