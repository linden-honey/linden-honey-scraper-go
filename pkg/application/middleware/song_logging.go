package middleware

import (
	"context"
	"log/slog"

	apisong "github.com/linden-honey/linden-honey-api-go/pkg/song"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain/song"
	"github.com/linden-honey/linden-honey-sdk-go/middleware"
)

// SongLoggingMiddleware returns a new instance of [middleware.Middleware[song.Service]] with top-level logging.
func SongLoggingMiddleware() middleware.Middleware[song.Service] {
	return func(next song.Service) song.Service {
		return &songLoggingMiddleware{
			next: next,
		}
	}
}

type songLoggingMiddleware struct {
	next song.Service
}

// GetSongs wraps the [song.Service] call with logging attached.
func (mw *songLoggingMiddleware) GetSongs(ctx context.Context) (out []apisong.Song, err error) {
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
func (mw *songLoggingMiddleware) GetSongsByScraperID(ctx context.Context, scrID string) (out []apisong.Song, err error) {
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
