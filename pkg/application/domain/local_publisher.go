package domain

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalPublisherService struct {
	fileName string
}

func NewLocalPublisherService(fileName string) *LocalPublisherService {
	return &LocalPublisherService{
		fileName: fileName,
	}
}

func (svc *LocalPublisherService) Publish(_ context.Context, in io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(svc.fileName), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output path: %w", err)
	}

	f, err := os.Create(svc.fileName)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	return nil
}
