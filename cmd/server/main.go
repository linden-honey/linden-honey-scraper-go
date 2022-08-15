package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/text/encoding/charmap"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/linden-honey/linden-honey-sdk-go/health"
	"github.com/linden-honey/linden-honey-sdk-go/validation"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/fetcher"
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
		defer func() {
			cancel()
			time.Sleep(3 * time.Second)
		}()
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	_ = logger.Log("msg", "initialization of the application")

	_ = logger.Log("msg", "initialize configuration")

	var cfg *config.Config
	{
		var err error
		if cfg, err = config.NewConfig(); err != nil {
			fatal(logger, fmt.Errorf("failed to initialize a config: %w", err))
		}
	}

	_ = logger.Log("msg", "initialize services")

	var scraperSvc scraper.Service
	{
		ss := make([]scraper.Service, 0)
		for id, scrCfg := range cfg.Scrapers {
			u, err := url.Parse(scrCfg.BaseURL)
			if err != nil {
				fatal(logger, fmt.Errorf("failed to parse scraper base url: %w", err))
			}

			f, err := fetcher.NewFetcher(
				fetcher.Config{
					BaseURL:        u,
					SourceEncoding: charmap.Windows1251,
				},
				fetcher.WithRetry(fetcher.RetryConfig{
					Retries:           5,
					Factor:            3,
					MinTimeout:        time.Second * 1,
					MaxTimeout:        time.Second * 6,
					MaxJitterInterval: time.Second,
				}),
			)
			if err != nil {
				fatal(logger, fmt.Errorf("failed to initialize a fetcher: %w", err))
			}

			p, err := parser.NewParser(id)
			if err != nil {
				fatal(logger, fmt.Errorf("failed to initialize a parser: %w", err))
			}

			v, err := validation.NewDelegate()
			if err != nil {
				fatal(logger, fmt.Errorf("failed to initialize a validator: %w", err))
			}

			scr, err := scraper.NewScraper(f, p, v)
			if err != nil {
				fatal(logger, fmt.Errorf("failed to initialize a scraper: %w", err))
			}

			s := scraper.LoggingMiddleware(
				log.With(
					logger,
					"component", "scraper", "scraper_id", id,
				),
			)(scr)

			ss = append(ss, s)
		}

		var err error
		scraperSvc, err = scraper.NewAggregator(ss...)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to initialize an aggregator: %w", err))
		}

		scraperSvc = scraper.LoggingMiddleware(
			log.With(
				logger,
				"component", "aggregator",
			),
		)(scraperSvc)
	}

	var docsSvc docs.Service
	{
		var err error
		docsSvc, err = docs.NewProvider("./api/openapi.json")
		if err != nil {
			fatal(logger, fmt.Errorf("failed to initialize docs provider: %w", err))
		}
	}

	_ = logger.Log("msg", "initialize endpoints")

	var scraperEndpoints scraper.Endpoints
	{
		scraperEndpoints = scraper.Endpoints{
			GetSong:     scraper.MakeGetSongEndpoint(scraperSvc),
			GetSongs:    scraper.MakeGetSongsEndpoint(scraperSvc),
			GetPreviews: scraper.MakeGetPreviewsEndpoint(scraperSvc),
		}
	}

	var docsEndpoints docs.Endpoints
	{
		docsEndpoints = docs.Endpoints{
			GetSpec: docs.MakeGetSpecEndpoint(docsSvc),
		}
	}

	_ = logger.Log("msg", "initialize http server")

	var httpServer *http.Server
	{
		router := mux.
			NewRouter().
			StrictSlash(true)

		logger := log.With(logger, "component", "http")

		if cfg.Health.Enabled {
			router.
				Path(cfg.Health.Path).
				Methods(http.MethodGet).
				Handler(
					health.NewHTTPHandler(
						health.MakeEndpoint(health.NewNopService()),
						logger,
					),
				)
		}

		// TODO: fix duplicate path prefixes

		router.PathPrefix("/api/songs").Handler(
			scraper.NewHTTPHandler(
				"/api/songs",
				scraperEndpoints,
				logger,
			),
		)
		router.PathPrefix("/").Handler(
			docs.NewHTTPHandler(
				"/",
				docsEndpoints,
				logger,
			),
		)

		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		httpServer = &http.Server{
			Addr:    addr,
			Handler: router,
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
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-sigc)
	}()

	_ = logger.Log("msg", "application started")
	_ = logger.Log("msg", "application stopped", "exit", <-errc)
}

func fatal(logger log.Logger, err error) {
	_ = logger.Log("err: %w", err)
	os.Exit(1)
}
