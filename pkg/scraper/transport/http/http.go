package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/endpoint"
)

func NewHTTPHandler(prefix string, endpoints *endpoint.Endpoints, logger log.Logger) http.Handler {
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	// initialize router
	r := mux.
		NewRouter().
		StrictSlash(true)

	// declare routes
	r.
		Path(prefix).
		Methods("GET").
		Queries("projection", "preview").
		Handler(httptransport.NewServer(
			endpoints.GetPreviews,
			decodeGetPreviewsRequest,
			encodeGetPreviewsResponse,
			opts...,
		))
	r.
		Path(prefix).
		Methods("GET").
		Handler(httptransport.NewServer(
			endpoints.GetSongs,
			decodeGetSongsRequest,
			encodeGetSongsResponse,
			opts...,
		))
	r.
		Path(fmt.Sprintf("%s/{id}", prefix)).
		Methods("GET").
		Handler(httptransport.NewServer(
			endpoints.GetSong,
			decodeGetSongRequest,
			encodeGetSongResponse,
			opts...,
		))

	return r
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

func decodeGetSongsRequest(_ context.Context, r *http.Request) (interface{}, error) {
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

func decodeGetPreviewsRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
