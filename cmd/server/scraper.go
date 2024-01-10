package main

import (
	"fmt"
	"net/url"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/parser"
)

func newScraper(cfg config.ScraperConfig, p scraper.Parser) (*scraper.Scraper, error) {
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scraper base url: %w", err)
	}

	f, err := fetcher.New(
		u,
		charmap.Windows1251,
		fetcher.WithRetry(&fetcher.RetryConfig{
			Attempts:    5,
			MinInterval: 2 * time.Second,
			MaxInterval: 10 * time.Second,
			Factor:      2 * time.Second,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a fetcher: %w", err)
	}

	return scraper.New(
		f,
		parser.NewGrobParser(),
		scraper.WithValidation(true),
	)
}
