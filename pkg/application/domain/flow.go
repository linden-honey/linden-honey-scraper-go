package domain

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

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

func (svc *FlowService) Run(ctx context.Context) error {
	var buf bytes.Buffer

	if err := svc.scrSvc.Scrape(ctx, &buf); err != nil {
		return fmt.Errorf("failed to scrape: %w", err)
	}

	if err := svc.pubSvc.Publish(ctx, &buf); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}
