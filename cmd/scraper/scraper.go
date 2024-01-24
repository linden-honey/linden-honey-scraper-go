package main

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/text/encoding/charmap"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/parser"
)

func newScrapers(cfg config.ScrapersConfig) (domain.Scrapers, error) {
	out := make(domain.Scrapers)

	{
		scr, err := newScraper(cfg.Grob, parser.NewGrobParser())
		if err != nil {
			return nil, fmt.Errorf("failed to init grob scraper: %w", err)
		}

		out["grob"] = scr
	}

	return out, nil
}

func newScraper(cfg config.ScraperConfig, p scraper.Parser) (*scraper.Scraper, error) {
	f, err := fetcher.New(
		&cfg.BaseURL,
		charmap.Windows1251, // TODO: use from cfg
		fetcherWithRetryOption(cfg.Retry),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a fetcher: %w", err)
	}

	return scraper.New(
		f,
		p,
		scraper.WithValidation(true), // TODO: use from cfg
	)
}

func fetcherWithRetryOption(cfg config.RetryConfig) fetcher.Option {
	if !cfg.Enabled {
		return fetcher.WithRetry(nil)
	}

	return fetcher.WithRetry(func(ctx context.Context, action func() (string, error)) (string, error) {
		return backoff.RetryWithData(
			action,
			backoff.WithContext(
				backoff.WithMaxRetries(
					&backoff.ExponentialBackOff{
						InitialInterval: cfg.MinInterval,
						MaxInterval:     cfg.MaxInterval,
					},
					uint64(cfg.Attempts),
				),
				ctx,
			),
		)
	})
}
