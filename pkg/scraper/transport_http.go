package scraper

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	sdkhttp "github.com/linden-honey/linden-honey-sdk-go/transport/http"
)

// NewHTTPHandler returns the new instance of http.Handler
func NewHTTPHandler(svc Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/{id}", makeGetSongHTTPHandlerFunc(svc))
	r.Get("/", makeGetSongsHTTPHandlerFunc(svc))
	r.Get("/previews", makeGetPreviewsHTTPHandlerFunc(svc))

	return r
}

func makeGetSongHTTPHandlerFunc(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		song, err := svc.GetSong(r.Context(), id)
		if err != nil {
			_ = sdkhttp.EncodeJSONError(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to get song by id=%s: %w", id, err),
			)

			return
		}

		_ = sdkhttp.EncodeJSONResponse(w, http.StatusOK, song)
	}
}

func makeGetSongsHTTPHandlerFunc(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songs, err := svc.GetSongs(r.Context())
		if err != nil {
			_ = sdkhttp.EncodeJSONError(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to get songs: %w", err),
			)

			return
		}

		_ = sdkhttp.EncodeJSONResponse(w, http.StatusOK, songs)
	}
}

func makeGetPreviewsHTTPHandlerFunc(svc Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		previews, err := svc.GetPreviews(r.Context())
		if err != nil {
			_ = sdkhttp.EncodeJSONError(
				w,
				http.StatusInternalServerError,
				fmt.Errorf("failed to get previews: %w", err),
			)

			return
		}

		_ = sdkhttp.EncodeJSONResponse(w, http.StatusOK, previews)
	}
}
