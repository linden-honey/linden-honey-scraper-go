package scraper

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

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

func (mw loggingMiddleware) GetSong(ctx context.Context, id string) (s *Song, err error) {
	_ = level.Debug(mw.logger).Log(
		"msg", "scrape a song",
		"id", id,
	)
	defer func() {
		if err == nil {
			_ = level.Debug(mw.logger).Log(
				"msg", "successfully scraped a song",
				"id", id,
				"title", s.Title,
			)
		} else {
			_ = level.Debug(mw.logger).Log(
				"msg", "failed to scrape a song",
				"id", id,
				"error", err,
			)
		}
	}()
	return mw.next.GetSong(ctx, id)
}

func (mw loggingMiddleware) GetSongs(ctx context.Context) (ss []Song, err error) {
	_ = level.Debug(mw.logger).Log(
		"msg", "start songs scraping",
	)
	defer func() {
		if err == nil {
			_ = level.Debug(mw.logger).Log(
				"msg", "songs scraping successfully finished",
				"count", len(ss),
			)
		} else {
			_ = level.Debug(mw.logger).Log(
				"msg", "songs scraping failed",
				"error", err,
			)
		}
	}()
	return mw.next.GetSongs(ctx)
}

func (mw loggingMiddleware) GetPreviews(ctx context.Context) (pp []Preview, err error) {
	_ = level.Debug(mw.logger).Log(
		"msg", "start previews scraping",
	)
	defer func() {
		if err == nil {
			_ = level.Debug(mw.logger).Log(
				"msg", "previews scraping successfully finished",
				"count", len(pp),
			)
		} else {
			_ = level.Debug(mw.logger).Log(
				"msg", "previews scraping failed",
				"error", err,
			)
		}
	}()
	return mw.next.GetPreviews(ctx)
}
