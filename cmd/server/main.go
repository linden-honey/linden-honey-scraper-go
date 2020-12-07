package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	docsendpoint "github.com/linden-honey/linden-honey-scraper-go/pkg/docs/endpoint"
	docssvc "github.com/linden-honey/linden-honey-scraper-go/pkg/docs/service"
	docshttptransport "github.com/linden-honey/linden-honey-scraper-go/pkg/docs/transport/http"
	scraperendpoint "github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/endpoint"
	scrapersvc "github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/service"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/service/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/service/parser"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/service/validator"
	scraperhttptransport "github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/transport/http"
)

func main() {
	// initialize logger
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	router := mux.
		NewRouter().
		StrictSlash(true)

	// TODO get URL from configuration
	u, _ := url.Parse("http://www.gr-oborona.ru")

	// initialize scrapper service
	scraperService := scrapersvc.NewService(
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
	scraperService = scrapersvc.LoggingMiddleware(
		log.With(
			logger,
			"component", "scraper",
			"scraper_source", u.String(),
		),
	)(scraperService)

	// initialize aggregator service
	scraperService = scrapersvc.NewAggregatorService(scraperService)
	scraperService = scrapersvc.LoggingMiddleware(
		log.With(
			logger,
			"component", "aggregator",
		),
	)(scraperService)

	// initialize scraper endpoints
	scraperEndpoints := scraperendpoint.NewEndpoints(scraperService)

	// initialize scraper http handler
	scraperHTTPHandler := scraperhttptransport.NewHTTPHandler("/api/songs", scraperEndpoints, logger)

	// register scraper handler
	router.PathPrefix("/api/songs").Handler(scraperHTTPHandler)

	// initialize docs service
	docsService := docssvc.NewService("./api/openapi-spec/openapi.json")

	// initialize docs endpoints
	docsEndpoints := docsendpoint.NewEndpoints(docsService)

	// initialize docs http handler
	docsHTTPHandler := docshttptransport.NewHTTPHandler("/", docsEndpoints, log.With(logger, "component", "http"))

	// register docs handler
	router.PathPrefix("/").Handler(docsHTTPHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		_ = level.Error(logger).Log("msg", "failed to serve http server", "err", err)
		os.Exit(1)
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
