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

	docsendpoint "github.com/linden-honey/linden-honey-scraper-go/pkg/docs/endpoint"
	docssvc "github.com/linden-honey/linden-honey-scraper-go/pkg/docs/service"
	docshttptransport "github.com/linden-honey/linden-honey-scraper-go/pkg/docs/transport/http"
	songendpoint "github.com/linden-honey/linden-honey-scraper-go/pkg/song/endpoint"
	songsvc "github.com/linden-honey/linden-honey-scraper-go/pkg/song/service"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/aggregator"
	songsvcmiddleware "github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/middleware"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper/parser"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper/validator"
	songhttptransport "github.com/linden-honey/linden-honey-scraper-go/pkg/song/transport/http"
)

func main() {
	// initialize logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	// initialize song service
	var songService songsvc.Service
	{
		// TODO get URL from configuration
		u, _ := url.Parse("http://www.gr-oborona.ru")

		// initialize scraper
		songService = scraper.NewScraper(
			fetcher.NewDefaultFetcherWithRetry(
				&fetcher.Properties{
					BaseURL:        u,
					SourceEncoding: charmap.Windows1251,
				},
				&fetcher.RetryProperties{
					Retries:    5,
					Factor:     3,
					MinTimeout: time.Second * 1,
					MaxTimeout: time.Second * 6,
				},
			),
			parser.NewDefaultParser(),
			validator.NewDefaultValidator(),
		)
		songService = songsvcmiddleware.LoggingMiddleware(
			log.With(
				logger,
				"component", "scraper",
				"source", u.String(),
			),
		)(songService)

		// initialize aggregator service
		songService = aggregator.NewAggregator(songService)
		songService = songsvcmiddleware.LoggingMiddleware(
			log.With(
				logger,
				"component", "aggregator",
			),
		)(songService)
	}

	// initialize songs endpoints
	var songEndpoints *songendpoint.Endpoints
	{
		songEndpoints = songendpoint.NewEndpoints(songService)
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
	var docsService docssvc.Service
	{
		docsService = docssvc.NewService("./api/openapi-spec/openapi.json")
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
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		_ = logger.Log("transport", "http", "addr", "0.0.0.0:8080")
		errs <- http.ListenAndServe(":8080", router)
	}()

	_ = logger.Log("exit", <-errs)
}
