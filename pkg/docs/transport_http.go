package docs

import (
	"context"
	"net/http"
	"path"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	swagger "github.com/swaggo/http-swagger"
)

// NewHTTPHandler returns the a instance of http.Handler
func NewHTTPHandler(prefix string, endpoints Endpoints, logger log.Logger) http.Handler {
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
			decodeGetSpecHTTPRequest,
			encodeGetSpecHTTPResponse,
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

func decodeGetSpecHTTPRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return GetSpecRequest{}, nil
}

func encodeGetSpecHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(GetSpecResponse)
	httptransport.SetContentType("application/json")(ctx, w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(*res.Spec))
	return nil
}
