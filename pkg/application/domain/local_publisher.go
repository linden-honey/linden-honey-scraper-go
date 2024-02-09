package domain

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type LocalPublisherService struct {
	logger *slog.Logger
}

func NewLocalPublisherService() *LocalPublisherService {
	return &LocalPublisherService{
		logger: slog.With(
			"component", "publisher",
			"publisher", "local",
		),
	}
}

func (svc *LocalPublisherService) Publish(ctx context.Context, artifactName string, in io.Reader) error {
	svc.logger.DebugContext(ctx, "ensuring directory structure exists")

	if err := os.MkdirAll(filepath.Dir(artifactName), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output path: %w", err)
	}

	svc.logger.DebugContext(ctx, "creating output file for the artifact")

	out, err := os.Create(artifactName)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	svc.logger.InfoContext(ctx, "writing input data to output file")

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy data to output file")
	}

	return nil
}
