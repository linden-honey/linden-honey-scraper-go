package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/controller"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/scraper"
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

	// Initialize song router
	songRouter := apiRouter.
		PathPrefix("/songs").
		Subrouter()

	//Declare song routes
	s := scraper.Create(&scraper.Properties{
		BaseURL: "http://www.gr-oborona.ru",
		Retry: scraper.RetryProperties{
			Retries:    5,
			Factor:     3,
			MinTimeout: time.Second * 1,
			MaxTimeout: time.Second * 6,
		},
	})
	songController := controller.SongController{
		Scraper: s,
	}
	songRouter.
		Path("/").
		Methods("GET").
		Queries("projection", "preview").
		HandlerFunc(songController.GetPreviews).
		Name("GetPreviews")
	songRouter.
		Path("/").
		Methods("GET").
		HandlerFunc(songController.GetSongs).
		Name("GetSongs")
	songRouter.
		Path("/{songId}").
		Methods("GET").
		HandlerFunc(songController.GetSong).
		Name("GetSong")

	log.Printf("Application is started on %d port!", 8080)
	log.Fatal(http.ListenAndServe(":8080", rootRouter))
}
