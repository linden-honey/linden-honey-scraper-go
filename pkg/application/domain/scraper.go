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

type ScraperService struct {
	scrapers Scrapers
	logger   *slog.Logger
}

type Scrapers map[string]Scraper

type Scraper interface {
	GetSongs(ctx context.Context) ([]song.Entity, error)
}

func NewScraperService(scrapers Scrapers, opts ...ScraperServiceOption) *ScraperService {
	svc := &ScraperService{
		scrapers: make(Scrapers),
		logger:   slog.With("component", "scraper"),
	}

	maps.Copy(svc.scrapers, scrapers)

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

type ScraperServiceOption func(*ScraperService)

// Scrape gets all the songs from multiple sources and writes them in json format to [io.Writer].
func (svc *ScraperService) Scrape(ctx context.Context, out io.Writer) error {
	songs, err := svc.getSongs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get songs: %w", err)
	}

	if err := json.NewEncoder(out).Encode(songs); err != nil {
		return fmt.Errorf("failed to encode songs as json: %w", err)
	}

	return nil
}

func (svc *ScraperService) getSongs(ctx context.Context) ([]song.Entity, error) {
	out := make([]song.Entity, 0)
	errs := make([]error, 0)
	for scrID, scr := range svc.scrapers {
		svc.logger.InfoContext(ctx, "getting songs", "scraper_id", scrID)

		songs, err := scr.GetSongs(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get songs from the scraper with id=%s: %w", scrID, err))
			continue
		}

		svc.logger.InfoContext(ctx, "songs successfully received", "scraper_id", scrID, "songs_count", len(songs))

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
