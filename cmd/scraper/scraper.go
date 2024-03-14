package main

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/parser"
)

func newScrapers(cfg config.ScrapersConfig) (map[string]domain.SongScraper, error) {
	out := make(map[string]domain.SongScraper)

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
		cfg.BaseURL,
		fetcher.WithEncoding(cfg.Encoding),
		fetcherWithRetry(cfg.Retry),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a fetcher: %w", err)
	}

	return scraper.New(
		f,
		p,
		scraper.WithValidation(cfg.Validation),
		scraper.WithSongResourcePath(cfg.SongResourcePath),
		scraper.WithSongMetadataListResourcePath(cfg.SongMetadataListResourcePath),
	)
}

func fetcherWithRetry(cfg config.RetryConfig) fetcher.Option {
	if !cfg.Enabled {
		return fetcher.WithRetry(nil)
	}

	return fetcher.WithRetry(func(ctx context.Context, action func() (string, error)) (string, error) {
		return backoff.RetryWithData(
			action,
			backoff.WithContext(
				backoff.WithMaxRetries(
					&backoff.ExponentialBackOff{
						InitialInterval:     cfg.MinInterval,
						RandomizationFactor: backoff.DefaultRandomizationFactor,
						Multiplier:          backoff.DefaultMultiplier,
						MaxInterval:         cfg.MaxInterval,
						MaxElapsedTime:      backoff.DefaultMaxElapsedTime,
						Stop:                backoff.Stop,
						Clock:               backoff.SystemClock,
					},
					uint64(cfg.Attempts),
				),
				ctx,
			),
		)
	})
}
