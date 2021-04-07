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

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/linden-honey/linden-honey-sdk-go/validation"

	"github.com/linden-honey/linden-honey-scraper-go/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs"
	docsendpoint "github.com/linden-honey/linden-honey-scraper-go/pkg/docs/endpoint"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs/service"
	docshttptransport "github.com/linden-honey/linden-honey-scraper-go/pkg/docs/transport/http"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
	songendpoint "github.com/linden-honey/linden-honey-scraper-go/pkg/song/endpoint"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/aggregator"
	songsvcmiddleware "github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/middleware"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper/parser"
	songhttptransport "github.com/linden-honey/linden-honey-scraper-go/pkg/song/transport/http"
)

func fatal(logger log.Logger, err error) {
	_ = logger.Log("err: %w", err)
	os.Exit(1)
}

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
	var songSvc song.Service
	{
		ss := make([]song.Service, 0)
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

			s := song.Compose(
				songsvcmiddleware.LoggingMiddleware(
					log.With(
						logger,
						"component", "scraper", "scraper_id", id,
					),
				),
			)(scr)

			ss = append(ss, s)
		}

		// initialize aggregator
		var err error
		songSvc, err = aggregator.NewAggregator(ss...)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to initialize an aggregator: %w", err))
		}

		songSvc = song.Compose(
			songsvcmiddleware.LoggingMiddleware(
				log.With(
					logger,
					"component", "aggregator",
				),
			),
		)(songSvc)
	}

	// initialize songs endpoints
	var songEndpoints songendpoint.Endpoints
	{
		songEndpoints = songendpoint.Endpoints{
			GetSong:     songendpoint.MakeGetSongEndpoint(songSvc),
			GetSongs:    songendpoint.MakeGetSongsEndpoint(songSvc),
			GetPreviews: songendpoint.MakeGetPreviewsEndpoint(songSvc),
		}
	}

	// initialize song http handler
	var songHTTPHandler http.Handler
	{
		songHTTPHandler = songhttptransport.NewHTTPHandler(
			"/api/songs",
			songEndpoints,
			logger,
		)
	}

	// initialize docs service
	var docsSvc docs.Service
	{
		var err error
		docsSvc, err = service.NewProvider("./api/openapi-spec/openapi.json")
		if err != nil {
			fatal(logger, fmt.Errorf("failed to initialize docs provider: %w", err))
		}
	}

	// initialize docs endpoints
	var docsEndpoints docsendpoint.Endpoints
	{
		docsEndpoints = docsendpoint.Endpoints{
			GetSpec: docsendpoint.MakeGetSpecEndpoint(docsSvc),
		}
	}

	// initialize docs http handler
	var docsHTTPHandler http.Handler
	{
		docsHTTPHandler = docshttptransport.NewHTTPHandler(
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
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := http.ListenAndServe(addr, httpHandler); err != nil {
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
