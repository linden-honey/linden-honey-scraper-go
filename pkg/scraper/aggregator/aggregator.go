package aggregator

import (
	"context"

	"github.com/linden-honey/linden-honey-go/pkg/song"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper"
)

// Aggregator represents the aggregation service implementation
type Aggregator struct {
	services []scraper.Service
}

// NewAggregator returns a pointer to the new instance of Aggregator or an error
func NewAggregator(services ...scraper.Service) (*Aggregator, error) {
	return &Aggregator{
		services: services,
	}, nil
}

// GetSong returns a pointer to the song or an error from aggregated services
func (a *Aggregator) GetSong(ctx context.Context, id string) (*song.Song, error) {
	errs := make([]error, 0)
	for _, svc := range a.services {
		s, err := svc.GetSong(ctx, id)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		return s, nil
	}

	return nil, NewAggregationError("failed to scrape a song", errs...)
}

// GetSongs returns songs or an error from aggregated services
func (a *Aggregator) GetSongs(ctx context.Context) ([]song.Song, error) {
	out := make([]song.Song, 0)
	errs := make([]error, 0)
	for _, svc := range a.services {
		songs, err := svc.GetSongs(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		out = append(out, songs...)
	}

	if len(errs) != 0 {
		return nil, NewAggregationError("failed to aggregate scraped songs", errs...)
	}

	return out, nil
}

// GetPreviews returns previews or an error from aggregated services
func (a *Aggregator) GetPreviews(ctx context.Context) ([]song.Meta, error) {
	out := make([]song.Meta, 0)
	errs := make([]error, 0)
	for _, svc := range a.services {
		previews, err := svc.GetPreviews(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		out = append(out, previews...)
	}

	if len(errs) != 0 {
		return nil, NewAggregationError("failed to aggregate scraped previews", errs...)
	}

	return out, nil
}