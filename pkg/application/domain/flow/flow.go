package flow

import (
	"context"
)

// TODO: think about structure
//  1. multiple packages
//     flow/simple/simple.go
//     flow/simple.go
//  2. single package - multiple interfaces
//     flow/flow.go (SimpleService, ComplexService)
//  3. other
//     maybe, different methods
type Service interface {
	Run(ctx context.Context) error
}
