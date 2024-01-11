package main

import (
	"context"
	"fmt"
	"net/url"

	"golang.org/x/text/encoding/charmap"

	"github.com/cenkalti/backoff/v4"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/fetcher"
)

func newScraper(cfg config.ScraperConfig, p scraper.Parser) (*scraper.Scraper, error) {
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scraper base url: %w", err)
	}

	f, err := fetcher.New(
		u,
		charmap.Windows1251,
		fetcher.WithRetry(func(ctx context.Context, action func() (string, error)) (string, error) {
			return backoff.RetryWithData(
				action,
				backoff.WithContext(
					backoff.WithMaxRetries(
						&backoff.ExponentialBackOff{
							InitialInterval: cfg.Retry.MinInterval,
							MaxInterval:     cfg.Retry.MaxInterval,
						},
						uint64(cfg.Retry.Attempts),
					),
					ctx,
				),
			)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a fetcher: %w", err)
	}

	return scraper.New(
		f,
		p,
		scraper.WithValidation(true),
	)
}
