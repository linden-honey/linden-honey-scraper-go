package domain

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain/flow"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain/publisher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain/scraper"
)

type FlowService struct {
	scrSvc scraper.Service
	pubSvc publisher.Service
	logger *slog.Logger
}

func NewFlowService(
	scrSvc scraper.Service,
	pubSvc publisher.Service,
) *FlowService {
	return &FlowService{
		scrSvc: scrSvc,
		pubSvc: pubSvc,
		logger: slog.With("component", "flow"),
	}
}

func (svc *FlowService) RunSimpleFlow(ctx context.Context, in flow.SimpleFlowInput) error {
	if err := in.Validate(); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	var buf bytes.Buffer

	svc.logger.Info("scraping songs")

	if err := svc.scrSvc.Scrape(ctx, &buf); err != nil {
		return fmt.Errorf("failed to scrape: %w", err)
	}

	svc.logger.Info("publishing result", "output", in.OutputFileName)

	if err := svc.pubSvc.Publish(ctx, in.OutputFileName, &buf); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}
