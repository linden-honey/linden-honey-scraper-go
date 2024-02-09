package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/linden-honey/linden-honey-scraper-go/cmd"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain/flow"
)

func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	{
		ctx = context.Background()
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
	}

	cmd.InitLogger()

	slog.Info("initialization of the application")

	slog.Info("initializing configuration")

	var cfg *config.Config
	{
		var err error
		cfg, err = config.New()
		if err != nil {
			cmd.Fatal(fmt.Errorf("failed to initialize a config: %w", err))
		}
	}

	slog.Info("initializing services")

	// TODO: think about initialization of services (single variable or multiple)
	var svc flow.Service
	{
		scrapers, err := newScrapers(cfg.Scrapers)
		if err != nil {
			cmd.Fatal(fmt.Errorf("failed to init scrapers: %w", err))
		}

		svc = domain.NewFlowService(
			domain.NewSongsScraperService(scrapers),
			domain.NewLocalPublisherService(),
		)
	}

	{

		if err := svc.RunSimpleFlow(ctx, flow.RunSimpleFlowRequest{
			ArtifactName: cfg.Output.FileName,
		}); err != nil {
			cmd.Fatal(fmt.Errorf("failed to run flow: %w", err))
		}
	}

	slog.Info("successfully finished")
}
