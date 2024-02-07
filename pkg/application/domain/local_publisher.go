package domain

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalPublisherService struct {
}

func NewLocalPublisherService() *LocalPublisherService {
	return &LocalPublisherService{}
}

func (svc *LocalPublisherService) Publish(_ context.Context, name string, in io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(name), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output path: %w", err)
	}

	out, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy data to output file")
	}

	return nil
}
