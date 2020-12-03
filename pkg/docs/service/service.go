package service

import (
	"context"
	"fmt"
	"io/ioutil"
)

// Service represents a specification provider interface
type Service interface {
	GetSpec(ctx context.Context) (string, error)
}

type specProvider struct {
	spec     string
	specPath string
}

func NewService(specPath string) Service {
	return &specProvider{
		specPath: specPath,
	}
}

func (sp *specProvider) GetSpec(ctx context.Context) (string, error) {
	if sp.spec == "" {
		if spec, err := ioutil.ReadFile(sp.specPath); err != nil {
			return "", fmt.Errorf("failed to read spec file: %w", err)
		} else {
			sp.spec = string(spec)
		}
	}
	return sp.spec, nil
}
