package domain

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sort"

	"github.com/linden-honey/linden-honey-api-go/pkg/song"
)

type ScraperService struct {
	scrapers Scrapers
}

type Scrapers map[string]Scraper

type Scraper interface {
	GetSongs(ctx context.Context) ([]song.Song, error)
}

func NewScraperService(scrapers Scrapers, opts ...ScraperServiceOption) *ScraperService {
	svc := &ScraperService{
		scrapers: make(Scrapers),
	}

	maps.Copy(svc.scrapers, scrapers)

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

type ScraperServiceOption func(*ScraperService)

// GetSongs scrapes all songs from multiple sources and returns a slice of [song.Song] instances or an error.
func (svc *ScraperService) GetSongs(ctx context.Context) ([]song.Song, error) {
	out := make([]song.Song, 0)
	errs := make([]error, 0)
	for scrID, scr := range svc.scrapers {
		ss, err := scr.GetSongs(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get songs from scraper with id=%s: %w", scrID, err))
			continue
		}

		out = append(out, ss...)
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Title < out[j].Title
	})

	return out, nil
}

// GetSongsByScraperID scrapes songs from the source returns a slice of [song.Song] instances or an error.
func (svc *ScraperService) GetSongsByScraperID(ctx context.Context, scrID string) ([]song.Song, error) {
	scr, ok := svc.scrapers[scrID]
	if !ok {
		return nil, fmt.Errorf("failed to resolve scraper by id=%s", scrID)
	}

	out, err := scr.GetSongs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get songs from scraper with id=%s: %w", scrID, err)
	}

	return out, nil
}
