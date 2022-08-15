package docs

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints represents endpoints definition
type Endpoints struct {
	GetSpec endpoint.Endpoint
}

func MakeGetSpecEndpoint(svc Service) endpoint.Endpoint {
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
