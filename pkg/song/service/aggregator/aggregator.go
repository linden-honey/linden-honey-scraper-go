package aggregator

import (
	"context"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/domain"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service"
)

// Aggregator represents the aggregation service implementation
type Aggregator struct {
	services []service.Service
}

// NewAggregator returns a pointer to the new instance of Aggregator
func NewAggregator(services ...service.Service) *Aggregator {
	return &Aggregator{
		services: services,
	}
}

func (a Aggregator) GetSong(ctx context.Context, id string) (*domain.Song, error) {
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

func (a Aggregator) GetSongs(ctx context.Context) ([]domain.Song, error) {
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

func (a Aggregator) GetPreviews(ctx context.Context) ([]domain.Preview, error) {
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
