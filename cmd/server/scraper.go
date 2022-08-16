package main

import (
	"fmt"
	"net/url"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/parser"
)

func newScraper(cfg config.ScraperConfig, p scraper.Parser) (*scraper.Scraper, error) {
	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scraper base url: %w", err)
	}

	f, err := fetcher.NewFetcher(
		baseURL,
		charmap.Windows1251,
		fetcher.WithRetry(fetcher.RetryConfig{
			Retries:           5,
			Factor:            3,
			MinTimeout:        time.Second * 1,
			MaxTimeout:        time.Second * 6,
			MaxJitterInterval: time.Second,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a fetcher: %w", err)
	}

	return scraper.NewScraper(
		f,
		parser.NewGrobParser(),
		scraper.WithValidation(true),
	)
}
