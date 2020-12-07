package service

import (
	"context"
	"fmt"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/domain"
)

type aggregationErr struct {
	msg     string
	reasons []error
}

func newAggregationErr(msg string, reasons ...error) *aggregationErr {
	return &aggregationErr{
		msg:     msg,
		reasons: reasons,
	}
}

func (err *aggregationErr) Error() string {
	return fmt.Sprintf("%s: %v", err.msg, err.reasons)
}

type aggregator struct {
	services []Service
}

func newAggregator(ss ...Service) *aggregator {
	return &aggregator{
		services: ss,
	}
}

func (a aggregator) GetSong(ctx context.Context, id string) (*domain.Song, error) {
	errs := make([]error, 0)
	for _, svc := range a.services {
		s, err := svc.GetSong(ctx, id)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		return s, nil
	}
	return nil, newAggregationErr("failed to scrape a song", errs...)
}

func (a aggregator) GetSongs(ctx context.Context) ([]domain.Song, error) {
	res := make([]domain.Song, 0)
	errs := make([]error, 0)
	for _, svc := range a.services {
		ss, err := svc.GetSongs(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		res = append(res, ss...)
	}
	if len(errs) != 0 {
		return nil, newAggregationErr("failed to aggregate scraped songs", errs...)
	}
	return res, nil
}

func (a aggregator) GetPreviews(ctx context.Context) ([]domain.Preview, error) {
	res := make([]domain.Preview, 0)
	errs := make([]error, 0)
	for _, svc := range a.services {
		pp, err := svc.GetPreviews(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		res = append(res, pp...)
	}
	if len(errs) != 0 {
		return nil, newAggregationErr("failed to aggregate scraped previews", errs...)
	}
	return res, nil
}
