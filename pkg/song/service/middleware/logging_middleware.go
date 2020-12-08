package middleware

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/domain"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next service.Service) service.Service {
		return loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   service.Service
}

func (mw loggingMiddleware) GetSong(ctx context.Context, id string) (s *domain.Song, err error) {
	_ = level.Debug(mw.logger).Log(
		"msg", "scrape a song",
		"song_id", id,
	)
	defer func() {
		if err == nil {
			_ = level.Debug(mw.logger).Log(
				"msg", "successfully scraped a song",
				"song_id", id,
				"song_title", s.Title,
			)
		} else {
			_ = level.Debug(mw.logger).Log(
				"msg", "failed to scrape a song",
				"song_id", id,
				"err", err,
			)
		}
	}()
	return mw.next.GetSong(ctx, id)
}

func (mw loggingMiddleware) GetSongs(ctx context.Context) (ss []domain.Song, err error) {
	_ = level.Debug(mw.logger).Log(
		"msg", "start songs scraping",
	)
	defer func() {
		if err == nil {
			_ = level.Debug(mw.logger).Log(
				"msg", "songs scraping successfully finished",
				"res_count", len(ss),
			)
		} else {
			_ = level.Debug(mw.logger).Log(
				"msg", "songs scraping failed",
				"err", err,
			)
		}
	}()
	return mw.next.GetSongs(ctx)
}

func (mw loggingMiddleware) GetPreviews(ctx context.Context) (pp []domain.Preview, err error) {
	_ = level.Debug(mw.logger).Log(
		"msg", "start previews scraping",
	)
	defer func() {
		if err == nil {
			_ = level.Debug(mw.logger).Log(
				"msg", "previews scraping successfully finished",
				"res_count", len(pp),
			)
		} else {
			_ = level.Debug(mw.logger).Log(
				"msg", "previews scraping failed",
				"err", err,
			)
		}
	}()
	return mw.next.GetPreviews(ctx)
}
