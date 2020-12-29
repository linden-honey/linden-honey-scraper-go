package main

import (
	"fmt"
	"net"
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
	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs/provider"
	docshttptransport "github.com/linden-honey/linden-honey-scraper-go/pkg/docs/transport/http"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/aggregator"
	songendpoint "github.com/linden-honey/linden-honey-scraper-go/pkg/song/endpoint"
	songmiddleware "github.com/linden-honey/linden-honey-scraper-go/pkg/song/middleware"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/scraper/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/scraper/parser"
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
			fatal(logger, "failed to initialize config", err)
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
				fatal(logger, "failed to initialize fetcher", err)
			}

			p, err := parser.NewParser(id)
			if err != nil {
				fatal(logger, "failed to initialize parser", err)
			}

			v, err := validation.NewDelegate()
			if err != nil {
				fatal(logger, "failed to initialize validator", err)
			}

			// initialize scraper
			scr, err := scraper.NewScraper(f, p, v)
			if err != nil {
				fatal(logger, "failed to initialize scraper", err)
			}

			s := songmiddleware.LoggingMiddleware(
				log.With(
					logger,
					"component", "scraper", "scraper_id", id,
				),
			)(scr)

			ss = append(ss, s)
		}

		// initialize aggregator
		var err error
		songService, err = aggregator.NewAggregator(ss...)
		if err != nil {
			fatal(logger, "failed to initialize aggregator", err)
		}

		songService = songmiddleware.LoggingMiddleware(
			log.With(
				logger,
				"component", "aggregator",
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
			fatal(logger, "failed to initialize docs service", err)
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

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		addr, err := net.ResolveTCPAddr(
			"",
			fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		)
		if err != nil {
			fatal(logger, "failed to resolve addr", err)
		}
		_ = logger.Log("transport", "http", "addr", addr.String())

		errs <- http.ListenAndServe(addr.String(), router)
	}()

	_ = logger.Log("exit", <-errs)
}
