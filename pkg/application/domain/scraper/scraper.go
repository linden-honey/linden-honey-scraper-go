package scraper

import (
	"context"
	"io"
)

// Service is an interface of scraper use-cases.
type Service interface {
	Scrape(ctx context.Context, out io.Writer) error
}
