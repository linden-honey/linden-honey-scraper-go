package aggregator

import (
	"context"
	"fmt"
	"sort"

	"github.com/linden-honey/linden-honey-api-go/pkg/song"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper"
)

// Aggregator is an implementation of the [scraper.Service]
// that aggregates results from multiple source.
type Aggregator struct {
	services []scraper.Service
}

// New returns a pointer to the new instance of [Aggregator] or an error.
func New(services ...scraper.Service) (*Aggregator, error) {
	return &Aggregator{
		services: services,
	}, nil
}

// GetSong tries to scrape a song by id from multiple services
// and returns a pointer to the new instance of [song.Song] or an error.
func (a *Aggregator) GetSong(ctx context.Context, id string) (*song.Song, error) {
	errs := make([]error, 0)
	for i, svc := range a.services {
		s, err := svc.GetSong(ctx, id)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get the song from services[%d]: %w", i, err))
			continue
		}

		return s, nil
	}

	return nil, NewAggregationError("failed to get the song from any service", errs...)
}

// GetSongs scrapes and aggregates all songs from multiple services
// and returns a slice of [song.Song] instances or an error.
func (a *Aggregator) GetSongs(ctx context.Context) ([]song.Song, error) {
	res := make([]song.Song, 0)
	errs := make([]error, 0)
	for i, svc := range a.services {
		ss, err := svc.GetSongs(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get songs from services[%d]: %w", i, err))
			continue
		}

		res = append(res, ss...)
	}

	if len(errs) != 0 {
		return nil, NewAggregationError("failed to aggregate songs", errs...)
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Title < res[j].Title
	})

	return res, nil
}

// GetPreviews scrapes songs metadata from multiple services
// and returns a slice of [song.Metadata] instances or an error.
func (a *Aggregator) GetPreviews(ctx context.Context) ([]song.Metadata, error) {
	res := make([]song.Metadata, 0)
	errs := make([]error, 0)
	for i, svc := range a.services {
		ps, err := svc.GetPreviews(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get previews from services[%d]: %w", i, err))
			continue
		}

		res = append(res, ps...)
	}

	if len(errs) != 0 {
		return nil, NewAggregationError("failed to aggregate previews", errs...)
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Title < res[j].Title
	})

	return res, nil
}
