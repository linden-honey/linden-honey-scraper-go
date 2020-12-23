package docs

import (
	"context"
)

// Service represents a docs service interface
type Service interface {
	GetSpec(ctx context.Context) (string, error)
}
