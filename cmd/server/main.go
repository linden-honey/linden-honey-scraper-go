package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	swagger "github.com/swaggo/http-swagger"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/controller"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/scraper"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/util/io"
)

func main() {
	// Initialize root router
	rootRouter := mux.
		NewRouter().
		StrictSlash(true)

	// Initialize api router
	apiRouter := rootRouter.
		PathPrefix("/api").
		Subrouter()

	// Initialize song controller
	url, _ := url.Parse("http://www.gr-oborona.ru")
	s := scraper.Create(&scraper.Properties{
		BaseURL: url,
		Retry: scraper.RetryProperties{
			Retries:    5,
			Factor:     3,
			MinTimeout: time.Second * 1,
			MaxTimeout: time.Second * 6,
		},
	})
	songController := &controller.SongController{
		Scraper: s,
	}

	// Initialize song router
	songRouter := apiRouter.
		PathPrefix("/songs").
		Subrouter()

	// Declare song routes
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

	// Initialize docs controller
	docsController := &controller.DocsController{
		Spec: io.MustReadContent("api/openapi-spec/openapi.json"),
	}

	// Initialize docs router
	docsRouter := rootRouter.
		PathPrefix("/").
		Subrouter()

	//Declare docs routes
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

	log.Printf("Application is started on %d port!", 8080)
	log.Fatal(http.ListenAndServe(":8080", rootRouter))
}
