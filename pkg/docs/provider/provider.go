package provider

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs"
)

// provider represents docs service implementation
type provider struct {
	spec     string
	specPath string
}

// NewService returns a pointer to the new instance of provider
func NewService(specPath string) docs.Service {
	return &provider{
		specPath: specPath,
	}
}

func (sp *provider) GetSpec(_ context.Context) (string, error) {
	if sp.spec == "" {
		if spec, err := ioutil.ReadFile(sp.specPath); err != nil {
			return "", fmt.Errorf("failed to read spec file: %w", err)
		} else {
			sp.spec = string(spec)
		}
	}
	return sp.spec, nil
}
