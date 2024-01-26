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
	"github.com/linden-honey/linden-honey-sdk-go/middleware"
)

type ScraperService struct {
	scrapers Scrapers
}

type Scrapers map[string]Scraper

type Scraper interface {
	GetSongs(ctx context.Context) ([]song.Entity, error)
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

// ScrapeSongs gets all the songs from multiple sources and writes them in json format to [io.Writer].
func (svc *ScraperService) ScrapeSongs(ctx context.Context, w io.Writer) error {
	songs, err := svc.getSongs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get songs: %w", err)
	}

	if err := json.NewEncoder(w).Encode(songs); err != nil {
		return fmt.Errorf("failed to encode songs as json: %w", err)
	}

	return nil
}

func (svc *ScraperService) getSongs(ctx context.Context) ([]song.Entity, error) {
	out := make([]song.Entity, 0)
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

// ScraperLoggingMiddleware returns a new instance of [middleware.Middleware[Scraper]] with top-level logging.
func ScraperLoggingMiddleware(scrID string) middleware.Middleware[Scraper] {
	return func(next Scraper) Scraper {
		return &scraperLoggingMiddleware{
			logger: slog.With("scraped_id", scrID),
			next:   next,
		}
	}
}

type scraperLoggingMiddleware struct {
	logger *slog.Logger
	next   Scraper
}

// GetSongs wraps the [song.Service] call with logging attached.
func (mw *scraperLoggingMiddleware) GetSongs(ctx context.Context) (out []song.Entity, err error) {
	mw.logger.InfoContext(ctx, "getting songs")

	defer func() {
		if err != nil {
			slog.ErrorContext(ctx, "failed to get songs", "err", err.Error())
		} else {
			slog.InfoContext(ctx, "successfully got songs", "songs_count", len(out))
		}
	}()

	return mw.next.GetSongs(ctx)
}
