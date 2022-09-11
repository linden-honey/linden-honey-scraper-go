package docs

import (
	"context"
	"fmt"
	"os"
)

// Provider represents a specification provider
type Provider struct {
	spec     *Spec
	specPath string
}

// NewProvider returns a pointer to a new instance of the provider or an error
func NewProvider(specPath string) (*Provider, error) {
	return &Provider{
		specPath: specPath,
	}, nil
}

// GetSpec returns a specification from the file system or cache or an error
func (sp *Provider) GetSpec(_ context.Context) (*Spec, error) {
	if sp.spec == nil {
		specBytes, err := os.ReadFile(sp.specPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read a spec file: %w", err)
		}

		spec := Spec(specBytes)
		sp.spec = &spec
	}
	return sp.spec, nil
}
