package aggregator

import (
	"context"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
)

// aggregator represents the aggregation service implementation
type aggregator struct {
	services []song.Service
}

// NewService returns a pointer to the new instance of aggregator
func NewService(services ...song.Service) song.Service {
	return &aggregator{
		services: services,
	}
}

func (a aggregator) GetSong(ctx context.Context, id string) (*song.Song, error) {
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

func (a aggregator) GetSongs(ctx context.Context) ([]song.Song, error) {
	res := make([]song.Song, 0)
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

func (a aggregator) GetPreviews(ctx context.Context) ([]song.Preview, error) {
	res := make([]song.Preview, 0)
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
