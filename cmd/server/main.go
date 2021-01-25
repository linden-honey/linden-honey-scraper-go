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
	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs/service/provider"
	docshttptransport "github.com/linden-honey/linden-honey-scraper-go/pkg/docs/transport/http"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
	songendpoint "github.com/linden-honey/linden-honey-scraper-go/pkg/song/endpoint"
	songmiddleware "github.com/linden-honey/linden-honey-scraper-go/pkg/song/middleware"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/aggregator"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper/parser"
	songhttptransport "github.com/linden-honey/linden-honey-scraper-go/pkg/song/transport/http"
)

func fatal(logger log.Logger, prefix string, err error) {
	err = fmt.Errorf("%s: %w", prefix, err)
	_ = logger.Log("err", err)
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
			fatal(logger, "failed to initialize a config", err)
		}
	}

	// initialize song service
	var songService song.Service
	{
		ss := make([]song.Service, 0)
		for id, scrCfg := range cfg.Application.Scrapers {
			u, err := url.Parse(scrCfg.BaseURL)
			if err != nil {
				fatal(logger, "failed to parse scraper base url", err)
			}

			f, err := fetcher.NewFetcherWithRetry(
				fetcher.Config{
					BaseURL:        u,
					SourceEncoding: charmap.Windows1251,
				},
				fetcher.RetryConfig{
					Retries:           5,
					Factor:            3,
					MinTimeout:        time.Second * 1,
					MaxTimeout:        time.Second * 6,
					MaxJitterInterval: time.Second,
				},
			)
			if err != nil {
				fatal(logger, "failed to initialize a fetcher", err)
			}

			p, err := parser.NewParser(id)
			if err != nil {
				fatal(logger, "failed to initialize a parser", err)
			}

			v, err := validation.NewDelegate()
			if err != nil {
				fatal(logger, "failed to initialize a validator", err)
			}

			// initialize scraper
			scr, err := scraper.NewScraper(f, p, v)
			if err != nil {
				fatal(logger, "failed to initialize a scraper", err)
			}

			s := song.Compose(
				songmiddleware.LoggingMiddleware(
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
		songService, err = aggregator.NewAggregator(ss...)
		if err != nil {
			fatal(logger, "failed to initialize an aggregator", err)
		}

		songService = song.Compose(
			songmiddleware.LoggingMiddleware(
				log.With(
					logger,
					"component", "aggregator",
				),
			),
		)(songService)
	}

	// initialize songs endpoints
	var songEndpoints *songendpoint.Endpoints
	{
		var err error
		songEndpoints, err = songendpoint.NewEndpoints(songService)
		if err != nil {
			fatal(logger, "failed to initialize endpoints", err)
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
	var docsService docs.Service
	{
		var err error
		docsService, err = provider.NewProvider("./api/openapi-spec/openapi.json")
		if err != nil {
			fatal(logger, "failed to initialize docs provider", err)
		}
	}

	// initialize docs endpoints
	var docsEndpoints *docsendpoint.Endpoints
	{
		docsEndpoints = docsendpoint.NewEndpoints(docsService)
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
	var router *mux.Router
	{
		router = mux.NewRouter().StrictSlash(true)

		// register song http handler
		router.PathPrefix("/api/songs").Handler(songHTTPHandler)

		// register docs handler
		router.PathPrefix("/").Handler(docsHTTPHandler)
	}

	errc := make(chan error, 1)

	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		_ = logger.Log("msg", "server started", "transport", "http", "addr", addr)
		errc <- http.ListenAndServe(addr, router)
	}()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-sigc)
	}()

	_ = logger.Log("exit", <-errc)
}
