package scraper

import (
	"context"
	"io"
)

// TODO: methods naming in the extension perspective
// 	- ScrapeFromSource(ctx context.Context, src string, out io.Writer) error
//	- ScrapePreviews(ctx context.Context, out io.Writer) error
//	- ???

// Service is an interface of scraper use-cases.
type Service interface {
	Scrape(ctx context.Context, out io.Writer) error
}
