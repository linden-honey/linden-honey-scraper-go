package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs/service"
)

type Endpoints struct {
	GetSpec endpoint.Endpoint
}

func NewEndpoints(svc service.Service) *Endpoints {
	return &Endpoints{
		GetSpec: makeGetSpecEndpoint(svc),
	}
}

func makeGetSpecEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(GetSpecRequest)
		spec, err := svc.GetSpec(ctx)
		return GetSpecResponse{
			Spec: spec,
		}, nil
	}
}
