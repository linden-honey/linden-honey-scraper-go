package flow

import (
	"context"
)

// Service is an interface of flow use-cases.
type Service interface {
	RunSimpleFlow(ctx context.Context, in RunSimpleFlowRequest) error
}
