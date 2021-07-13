package docs

import (
	"context"
	"fmt"
	"io/ioutil"
)

// Provider represents spec provider
type Provider struct {
	spec     *Spec
	specPath string
}

// NewProvider returns a pointer to the new instance of Provider or an error
func NewProvider(specPath string) (*Provider, error) {
	return &Provider{
		specPath: specPath,
	}, nil
}

// GetSpec returns specification from fs or cache or an error
func (sp *Provider) GetSpec(_ context.Context) (*Spec, error) {
	if sp.spec == nil {
		specBytes, err := ioutil.ReadFile(sp.specPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read a spec file: %w", err)
		}

		spec := Spec(specBytes)
		sp.spec = &spec
	}
	return sp.spec, nil
}