package endpoint

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs"
)

// Endpoints represents endpoints definition
type Endpoints struct {
	GetSpec endpoint.Endpoint
}

// NewEndpoints returns a pointer to the new instance of Endpoints or an error
func NewEndpoints(svc docs.Service) *Endpoints {
	return &Endpoints{
		GetSpec: makeGetSpecEndpoint(svc),
	}
}

func makeGetSpecEndpoint(svc docs.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(GetSpecRequest)
		spec, err := svc.GetSpec(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get spec: %w", err)
		}

		return GetSpecResponse{
			Spec: spec,
		}, nil
	}
}
