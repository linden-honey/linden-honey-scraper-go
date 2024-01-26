package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/linden-honey/linden-honey-scraper-go/cmd"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain/scraper"
)

func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	{
		ctx = context.Background()
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
	}

	cmd.InitLogger()

	slog.Info("initialization of the application")

	slog.Info("initializing configuration")

	var cfg *config.Config
	{
		var err error
		cfg, err = config.New()
		if err != nil {
			cmd.Fatal(fmt.Errorf("failed to initialize a config: %w", err))
		}
	}

	slog.Info("initializing services")

	var svc scraper.Service
	{
		scrapers, err := newScrapers(cfg.Scrapers)
		if err != nil {
			cmd.Fatal(fmt.Errorf("failed to init scrapers: %w", err))
		}

		svc = domain.NewScraperService(scrapers)
	}

	{
		slog.Info("scraping songs")
		if err := svc.ScrapeSongs(ctx, os.Stdout); err != nil {
			cmd.Fatal(fmt.Errorf("failed to scrape songs: %w", err))
		}
	}

	slog.Info("successfully finished")
}