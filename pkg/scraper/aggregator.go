package scraper

import (
	"context"
	"fmt"

	"github.com/linden-honey/linden-honey-go/pkg/song"
)

// Aggregator represents the aggregation service implementation
type Aggregator struct {
	services []Service
}

// NewAggregator returns a pointer to the new instance of Aggregator or an error
func NewAggregator(services ...Service) (*Aggregator, error) {
	return &Aggregator{
		services: services,
	}, nil
}

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
	return fmt.Sprintf("%s: %s", err.msg, err.reasons)
}

// GetSong returns a pointer to the song or an error from aggregated services
func (a Aggregator) GetSong(ctx context.Context, id string) (*song.Song, error) {
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

// GetSongs returns songs or an error from aggregated services
func (a Aggregator) GetSongs(ctx context.Context) ([]song.Song, error) {
	out := make([]song.Song, 0)
	errs := make([]error, 0)
	for _, svc := range a.services {
		ss, err := svc.GetSongs(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		out = append(out, ss...)
	}

	if len(errs) != 0 {
		return nil, newAggregationErr("failed to aggregate scraped songs", errs...)
	}

	return out, nil
}

// GetPreviews returns previews or an error from aggregated services
func (a Aggregator) GetPreviews(ctx context.Context) ([]song.Preview, error) {
	out := make([]song.Preview, 0)
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
		return nil, newAggregationErr("failed to aggregate scraped previews", errs...)
	}

	return out, nil
}
