package scraper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

// NewHTTPHandler returns the new instance of http.Handler
func NewHTTPHandler(prefix string, endpoints Endpoints, logger log.Logger) http.Handler {
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	// initialize router
	router := mux.
		NewRouter().
		StrictSlash(true)

	// declare routes
	router.
		Path(path.Clean(prefix)).
		Methods("GET").
		Queries("view", "preview").
		Handler(httptransport.NewServer(
			endpoints.GetPreviews,
			decodeGetPreviewsHTTPRequest,
			encodeGetPreviewsHTTPResponse,
			opts...,
		))
	router.
		Path(path.Clean(prefix)).
		Methods("GET").
		Handler(httptransport.NewServer(
			endpoints.GetSongs,
			decodeGetSongsHTTPRequest,
			encodeGetSongsHTTPResponse,
			opts...,
		))
	router.
		Path(path.Join(prefix, "{id}")).
		Methods("GET").
		Handler(httptransport.NewServer(
			endpoints.GetSong,
			decodeGetSongHTTPRequest,
			encodeGetSongHTTPResponse,
			opts...,
		))

	return router
}

func decodeGetSongHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("missed song id")
	}
	return GetSongRequest{
		ID: id,
	}, nil
}

func encodeGetSongHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(GetSongResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Result); err != nil {
		return fmt.Errorf("failed to encode get song response: %w", err)
	}
	return nil
}

func decodeGetSongsHTTPRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return GetSongsRequest{}, nil
}

func encodeGetSongsHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(GetSongsResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Results); err != nil {
		return fmt.Errorf("failed to encode get songs response: %w", err)
	}

	return nil
}

func decodeGetPreviewsHTTPRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return GetPreviewsRequest{}, nil
}

func encodeGetPreviewsHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(GetPreviewsResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Results); err != nil {
		return fmt.Errorf("failed to encode get previews response: %w", err)
	}

	return nil
}
