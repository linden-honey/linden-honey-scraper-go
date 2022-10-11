package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-kit/log"

	"github.com/linden-honey/linden-honey-sdk-go/health"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/aggregator"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/parser"
)

func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	{
		ctx = context.Background()
		ctx, cancel = context.WithCancel(ctx)
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer func() {
			cancel()
		}()
	}

	var logger log.Logger
	{
		logger = newLogger()
	}

	_ = logger.Log("msg", "initialization of the application")

	_ = logger.Log("msg", "initialize configuration")

	var cfg *config.Config
	{
		var err error
		cfg, err = config.NewConfig()
		if err != nil {
			fatal(logger, fmt.Errorf("failed to initialize a config: %w", err))
		}
	}

	_ = logger.Log("msg", "initialize services")

	var scrSvc scraper.Service
	{
		var err error
		var grobScrSvc scraper.Service
		{
			grobScrSvc, err = newScraper(cfg.Scrapers.Grob, parser.NewGrobParser())
			if err != nil {
				fatal(logger, fmt.Errorf("failed to initialize grob scraper: %w", err))
			}

			grobScrSvc = scraper.LoggingMiddleware(
				log.With(
					logger,
					"component", "scraper",
					"scraper_id", "grob",
				),
			)(grobScrSvc)
		}

		scrSvc, err = aggregator.NewAggregator(
			grobScrSvc,
		)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to initialize an aggregator: %w", err))
		}

		scrSvc = scraper.LoggingMiddleware(
			log.With(logger, "component", "aggregator"),
		)(scrSvc)
	}

	_ = logger.Log("msg", "initialize http server")

	var httpServer *http.Server
	{
		r := chi.NewRouter()
		r.Use(middleware.Recoverer)

		if cfg.Health.Enabled {
			r.Handle(cfg.Health.Path, health.NewHTTPHandler(health.NewNopService()))
		}

		r.Route("/api", func(r chi.Router) {
			r.Mount("/songs", scraper.NewHTTPHandler(scrSvc))
		})

		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		httpServer = &http.Server{
			Addr:    addr,
			Handler: r,
		}

		defer func() {
			if err := httpServer.Shutdown(ctx); err != nil {
				fatal(logger, fmt.Errorf("failed to shutdown http server: %w", err))
			}
		}()
	}

	errc := make(chan error, 1)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errc <- err
		}
	}()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		errc <- fmt.Errorf("%s", <-sigc)
	}()

	_ = logger.Log("msg", "application started")
	_ = logger.Log("msg", "application stopped", "exit", <-errc)
}
