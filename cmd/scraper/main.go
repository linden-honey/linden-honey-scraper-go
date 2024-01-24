package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	sdkmiddleware "github.com/linden-honey/linden-honey-sdk-go/middleware"

	"github.com/linden-honey/linden-honey-scraper-go/cmd"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/middleware"
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
		svc = sdkmiddleware.Compose(
			middleware.ScraperLoggingMiddleware(),
		)(svc)
	}

	{
		// TODO: initialize some transport here (s3, localfs, etc)
		_ = svc
		_ = ctx
	}

	errc := make(chan error, 1)

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-sigc)
	}()

	slog.Info("application started")
	slog.Info("application stopped", "exit", <-errc)
}
