package http

import (
	"context"
	"errors"
	"net/http"
	"path"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/endpoint"
)

// NewHTTPHandler returns the new instance of http.Handler
func NewHTTPHandler(prefix string, endpoints endpoint.Endpoints, logger log.Logger) http.Handler {
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
			decodeGetPreviewsRequest,
			encodeGetPreviewsResponse,
			opts...,
		))
	router.
		Path(path.Clean(prefix)).
		Methods("GET").
		Handler(httptransport.NewServer(
			endpoints.GetSongs,
			decodeGetSongsRequest,
			encodeGetSongsResponse,
			opts...,
		))
	router.
		Path(path.Join(prefix, "{id}")).
		Methods("GET").
		Handler(httptransport.NewServer(
			endpoints.GetSong,
			decodeGetSongRequest,
			encodeGetSongResponse,
			opts...,
		))

	return router
}

func decodeGetSongRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("missed song id")
	}
	return endpoint.GetSongRequest{
		ID: id,
	}, nil
}

func encodeGetSongResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoint.GetSongResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Result); err != nil {
		return err
	}
	return nil
}

func decodeGetSongsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return endpoint.GetSongsRequest{}, nil
}

func encodeGetSongsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoint.GetSongsResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Results); err != nil {
		return err
	}
	return nil
}

func decodeGetPreviewsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return endpoint.GetPreviewsRequest{}, nil
}

func encodeGetPreviewsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoint.GetPreviewsResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	if err := httptransport.EncodeJSONResponse(ctx, w, res.Results); err != nil {
		return err
	}
	return nil
}
