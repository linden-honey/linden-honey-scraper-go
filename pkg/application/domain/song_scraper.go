package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"sort"

	"github.com/linden-honey/linden-honey-api-go/pkg/application/domain/song"
)

// SongScraperService is an implementation of [scraper.Service] for scraping songs.
type SongScraperService struct {
	scrapers map[string]SongScraper
	logger   *slog.Logger
}

// SongScraper is an API for scraping songs.
type SongScraper interface {
	GetSongs(ctx context.Context) ([]song.Entity, error)
}

// NewSongsScraperService returns a pointer to the new instance of [SongScraperService].
func NewSongsScraperService(scrapers map[string]SongScraper, opts ...SongScraperServiceOption) *SongScraperService {
	svc := &SongScraperService{
		scrapers: make(map[string]SongScraper),
		logger: slog.With(
			"component", "scraper",
			"scraper", "song",
		),
	}

	maps.Copy(svc.scrapers, scrapers)

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

// SongScraperServiceOption set optional parameters for the [SongScraperService].
type SongScraperServiceOption func(*SongScraperService)

// Scrape scrapes songs from multiple sources and writes the result in json format to [io.Writer]
func (svc *SongScraperService) Scrape(ctx context.Context, out io.Writer) error {
	svc.logger.InfoContext(ctx, "getting songs")

	songs, err := svc.getSongs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get songs: %w", err)
	}

	svc.logger.InfoContext(ctx, "encoding songs in json format", "songs_count", len(songs))

	if err := json.NewEncoder(out).Encode(songs); err != nil {
		return fmt.Errorf("failed to encode songs in json format: %w", err)
	}

	return nil
}

func (svc *SongScraperService) getSongs(ctx context.Context) ([]song.Entity, error) {
	out := make([]song.Entity, 0)
	errs := make([]error, 0)
	for scrID, scr := range svc.scrapers {
		svc.logger.DebugContext(ctx, "getting songs from the scraper", "scraper_id", scrID)

		songs, err := scr.GetSongs(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get songs from the scraper with id=%s: %w", scrID, err))
			continue
		}

		svc.logger.DebugContext(ctx, "songs successfully received from the scraper", "scraper_id", scrID, "songs_count", len(songs))

		out = append(out, songs...)
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	svc.logger.InfoContext(ctx, "all songs successfully received", "songs_count", len(out))

	sort.Slice(out, func(i, j int) bool {
		return out[i].Title < out[j].Title
	})

	return out, nil
}
