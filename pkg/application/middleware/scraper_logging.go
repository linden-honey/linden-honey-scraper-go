package middleware

import (
	"context"
	"log/slog"

	"github.com/linden-honey/linden-honey-api-go/pkg/song"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain/scraper"
	"github.com/linden-honey/linden-honey-sdk-go/middleware"
)

// ScraperServiceLoggingMiddleware returns a new instance of [middleware.Middleware[song.Service]] with top-level logging.
func ScraperServiceLoggingMiddleware() middleware.Middleware[scraper.Service] {
	return func(next scraper.Service) scraper.Service {
		return &scraperServiceLoggingMiddleware{
			next: next,
		}
	}
}

type scraperServiceLoggingMiddleware struct {
	next scraper.Service
}

// GetSongs wraps the [song.Service] call with logging attached.
func (mw *scraperServiceLoggingMiddleware) GetSongs(ctx context.Context) (out []song.Song, err error) {
	slog.Info("getting songs")

	defer func() {
		if err != nil {
			slog.ErrorContext(ctx, "failed to get songs", "err", err.Error())
		} else {
			slog.Info("successfully got songs", "songs_count", len(out))
		}
	}()

	return mw.next.GetSongs(ctx)
}

// GetSongsByScraperID wraps the [song.Service] call with logging attached.
func (mw *scraperServiceLoggingMiddleware) GetSongsByScraperID(ctx context.Context, scrID string) (out []song.Song, err error) {
	slog.InfoContext(ctx, "getting songs by source id", "songs_scraper_id", scrID)

	defer func() {
		if err != nil {
			slog.ErrorContext(ctx, "failed to get songs by source id", "songs_scraper_id", scrID, "err", err.Error())
		} else {
			slog.InfoContext(ctx, "successfully got songs by source id", "songs_count", len(out), "songs_source", scrID)
		}
	}()

	return mw.next.GetSongs(ctx)
}
