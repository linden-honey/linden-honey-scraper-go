package http

import (
	"context"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	swagger "github.com/swaggo/http-swagger"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs/endpoint"
)

// NewHTTPHandler returns the new instance of http.Handler
func NewHTTPHandler(prefix string, endpoints *endpoint.Endpoints, logger log.Logger) http.Handler {
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	r := mux.
		NewRouter().
		StrictSlash(true)

	specPath := path.Join(prefix, "docs")
	r.
		Path(specPath).
		Methods("GET").
		Handler(httptransport.NewServer(
			endpoints.GetSpec,
			decodeGetSpecRequest,
			encodeGetSpecResponse,
			opts...,
		))
	r.
		PathPrefix(path.Clean(prefix)).
		Methods("GET").
		Handler(swagger.Handler(
			swagger.URL(specPath),
		))

	return r
}

func decodeGetSpecRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return endpoint.GetSpecRequest{}, nil
}

func encodeGetSpecResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoint.GetSpecResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(*res.Spec))
	return nil
}
