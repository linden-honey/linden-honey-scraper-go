package scraper

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/linden-honey/linden-honey-go/pkg/song"
)

// LoggingMiddleware returns logging middleware for scraper service
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

// GetSong proxies call to service with logging
func (mw loggingMiddleware) GetSong(ctx context.Context, id string) (s *song.Song, err error) {
	_ = level.Debug(mw.logger).Log(
		"msg", "getting a song",
		"song_id", id,
	)
	defer func() {
		if err == nil {
			_ = level.Debug(mw.logger).Log(
				"msg", "successfully got a song",
				"song_id", id,
				"song_title", s.Title,
			)
		} else {
			_ = level.Debug(mw.logger).Log(
				"msg", "failed to get a song",
				"song_id", id,
				"err", err,
			)
		}
	}()
	return mw.next.GetSong(ctx, id)
}

// GetSongs proxies call to service with logging
func (mw loggingMiddleware) GetSongs(ctx context.Context) (ss []song.Song, err error) {
	_ = level.Debug(mw.logger).Log(
		"msg", "getting songs",
	)
	defer func() {
		if err == nil {
			_ = level.Debug(mw.logger).Log(
				"msg", "successfully got songs",
				"count", len(ss),
			)
		} else {
			_ = level.Debug(mw.logger).Log(
				"msg", "failed to get songs",
				"err", err,
			)
		}
	}()
	return mw.next.GetSongs(ctx)
}

// GetPreviews proxies call to service with logging
func (mw loggingMiddleware) GetPreviews(ctx context.Context) (pp []song.Meta, err error) {
	_ = level.Debug(mw.logger).Log(
		"msg", "getting previews",
	)
	defer func() {
		if err == nil {
			_ = level.Debug(mw.logger).Log(
				"msg", "successfully got previews",
				"count", len(pp),
			)
		} else {
			_ = level.Debug(mw.logger).Log(
				"msg", "failed to get previews",
				"err", err,
			)
		}
	}()
	return mw.next.GetPreviews(ctx)
}
