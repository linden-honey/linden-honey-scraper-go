package provider

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs"
)

// Provider represents spec provider
type Provider struct {
	spec     *docs.Spec
	specPath string
}

// NewService returns a pointer to the new instance of Provider or an error
func NewService(specPath string) (*Provider, error) {
	return &Provider{
		specPath: specPath,
	}, nil
}

// GetSpec returns specification from fs or cache or an error
func (sp *Provider) GetSpec(_ context.Context) (*docs.Spec, error) {
	if sp.spec == nil {
		specBytes, err := ioutil.ReadFile(sp.specPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read spec file: %w", err)
		}

		spec := docs.Spec(specBytes)
		sp.spec = &spec
	}
	return sp.spec, nil
}
