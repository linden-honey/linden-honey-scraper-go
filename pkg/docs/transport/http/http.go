package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	swagger "github.com/swaggo/http-swagger"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs/endpoint"
)

func NewHTTPHandler(prefix string, endpoints *endpoint.Endpoints, logger log.Logger) http.Handler {
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	r := mux.
		NewRouter().
		StrictSlash(true)

	r.
		Path(fmt.Sprintf("%s/docs", prefix)).
		Methods("GET").
		Handler(httptransport.NewServer(
			endpoints.GetSpec,
			decodeGetSpecRequest,
			encodeGetSpecResponse,
			opts...,
		))
	r.
		PathPrefix(prefix).
		Methods("GET").
		Handler(swagger.Handler(
			swagger.URL("/docs"),
		))

	return r
}

func decodeGetSpecRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoint.GetSpecRequest{}, nil
}

func encodeGetSpecResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoint.GetSpecResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(res.Spec))
	return nil
}
