package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain/scraper"
	sdkhttp "github.com/linden-honey/linden-honey-sdk-go/transport/http"
)

// NewScraperHandler returns a new instance of [http.Handler].
func NewScraperHandler(svc scraper.Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/", makeScraperGetSongsHandlerFunc(svc))
	r.Get("/{scraperID}", makeScraperGetSongsByScraperIDHandlerFunc(svc))

	return r
}

func makeScraperGetSongsHandlerFunc(svc scraper.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := svc.GetSongs(r.Context())
		if err != nil {
			_ = sdkhttp.EncodeJSONError(
				w,
				http.StatusUnprocessableEntity,
				fmt.Errorf("failed to get songs: %w", err),
			)

			return
		}

		_ = sdkhttp.EncodeJSONResponse(w, http.StatusOK, out)
	}
}

func makeScraperGetSongsByScraperIDHandlerFunc(svc scraper.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scrID := chi.URLParam(r, "scraperID")
		out, err := svc.GetSongsByScraperID(r.Context(), scrID)
		if err != nil {
			_ = sdkhttp.EncodeJSONError(
				w,
				http.StatusNotFound,
				fmt.Errorf("failed to get song by id=%s: %w", scrID, err),
			)

			return
		}

		_ = sdkhttp.EncodeJSONResponse(w, http.StatusOK, out)
	}
}
