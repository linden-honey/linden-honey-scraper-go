package main

import (
	"net/http"
	"net/url"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	swagger "github.com/swaggo/http-swagger"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/controller"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/parser"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/validator"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/util/io"
)

func main() {
	// initialize logger
	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// initialize root router
	rootRouter := mux.
		NewRouter().
		StrictSlash(true)

	// initialize api router
	apiRouter := rootRouter.
		PathPrefix("/api").
		Subrouter()

	//parse
	u, err := url.Parse("http://www.gr-oborona.ru")
	if err != nil {
		logger.Fatal("Can't parse base URL", err)
	}

	// initialize scrapper
	s := scraper.NewDefaultScraper(
		logger,
		fetcher.NewDefaultFetcherWithRetry(
			logger,
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
		parser.NewDefaultParser(logger),
		validator.NewDefaultValidator(logger),
	)

	// initialize song controller
	songController := controller.NewSongController(
		logger,
		s,
	)

	// initialize song router
	songRouter := apiRouter.
		PathPrefix("/songs").
		Subrouter()

	// declare song routes
	songRouter.
		Path("/").
		Methods("GET").
		Queries("projection", "preview").
		HandlerFunc(songController.GetPreviews).
		Name("getPreviews")
	songRouter.
		Path("/").
		Methods("GET").
		HandlerFunc(songController.GetSongs).
		Name("getSongs")
	songRouter.
		Path("/{songId}").
		Methods("GET").
		HandlerFunc(songController.GetSong).
		Name("getSong")

	// initialize docs controller
	spec := io.MustReadContent("api/openapi-spec/openapi.json")
	docsController := controller.NewDocsController(
		logger,
		spec,
	)

	// initialize docs router
	docsRouter := rootRouter.
		PathPrefix("/").
		Subrouter()

	// declare docs routes
	docsRouter.
		Path("/api-docs").
		Methods("GET").
		HandlerFunc(docsController.GetSpec).
		Name("getApiDocs")
	docsRouter.
		PathPrefix("/").
		Methods("GET").
		Handler(swagger.Handler(
			swagger.URL("/api-docs"),
		)).
		Name("swagger")

	logger.Printf("Application is started on %d port!", 8080)
	logger.Fatal(http.ListenAndServe(":8080", rootRouter))
}
