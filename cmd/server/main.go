package main

import (
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

	"github.com/linden-honey/linden-honey-sdk-go/validation"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/parser"
)

func main() {
	// initialize logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	// initialize config
	var cfg *config.Config
	{
		var err error
		if cfg, err = config.NewConfig(); err != nil {
			fatal(logger, fmt.Errorf("failed to initialize a config: %w", err))
		}
	}

	// initialize song service
	var songSvc scraper.Service
	{
		ss := make([]scraper.Service, 0)
		for id, scrCfg := range cfg.Application.Scrapers {
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

			// initialize scraper
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

		// initialize aggregator
		var err error
		songSvc, err = scraper.NewAggregator(ss...)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to initialize an aggregator: %w", err))
		}

		songSvc = scraper.LoggingMiddleware(
			log.With(
				logger,
				"component", "aggregator",
			),
		)(songSvc)
	}

	// initialize songs endpoints
	var songEndpoints scraper.Endpoints
	{
		songEndpoints = scraper.Endpoints{
			GetSong:     scraper.MakeGetSongEndpoint(songSvc),
			GetSongs:    scraper.MakeGetSongsEndpoint(songSvc),
			GetPreviews: scraper.MakeGetPreviewsEndpoint(songSvc),
		}
	}

	// initialize song http handler
	var songHTTPHandler http.Handler
	{
		songHTTPHandler = scraper.NewHTTPHandler(
			"/api/songs",
			songEndpoints,
			logger,
		)
	}

	// initialize docs service
	var docsSvc docs.Service
	{
		var err error
		docsSvc, err = docs.NewProvider("./api/openapi.json")
		if err != nil {
			fatal(logger, fmt.Errorf("failed to initialize docs provider: %w", err))
		}
	}

	// initialize docs endpoints
	var docsEndpoints docs.Endpoints
	{
		docsEndpoints = docs.Endpoints{
			GetSpec: docs.MakeGetSpecEndpoint(docsSvc),
		}
	}

	// initialize docs http handler
	var docsHTTPHandler http.Handler
	{
		docsHTTPHandler = docs.NewHTTPHandler(
			"/",
			docsEndpoints,
			log.With(logger, "component", "http"),
		)
	}

	// initialize router
	var httpHandler http.Handler
	{
		router := mux.NewRouter().StrictSlash(true)

		router.PathPrefix("/api/songs").Handler(songHTTPHandler)
		router.PathPrefix("/").Handler(docsHTTPHandler)

		httpHandler = router
	}

	errc := make(chan error, 1)

	go func() {
		if err := http.ListenAndServe(cfg.Server.Addr, httpHandler); err != nil {
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
