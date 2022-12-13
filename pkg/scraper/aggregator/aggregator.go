package aggregator

import (
	"context"

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

// GetSong returns a pointer to a song or an error from aggregated services.
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

// GetSongs returns songs or an error from aggregated services.
func (a *Aggregator) GetSongs(ctx context.Context) ([]song.Song, error) {
	res := make([]song.Song, 0)
	errs := make([]error, 0)
	for _, svc := range a.services {
		songs, err := svc.GetSongs(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		res = append(res, songs...)
	}

	if len(errs) != 0 {
		return nil, NewAggregationError("failed to aggregate scraped songs", errs...)
	}

	return res, nil
}

// GetPreviews returns previews or an error from aggregated services.
func (a *Aggregator) GetPreviews(ctx context.Context) ([]song.Metadata, error) {
	res := make([]song.Metadata, 0)
	errs := make([]error, 0)
	for _, svc := range a.services {
		previews, err := svc.GetPreviews(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		res = append(res, previews...)
	}

	if len(errs) != 0 {
		return nil, NewAggregationError("failed to aggregate scraped previews", errs...)
	}

	return res, nil
}
