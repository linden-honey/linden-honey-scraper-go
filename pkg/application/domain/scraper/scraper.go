package scraper

import (
	"context"

	"github.com/linden-honey/linden-honey-api-go/pkg/song"
)

// Service is an interface of song use-cases.
type Service interface {
	GetSongs(ctx context.Context) ([]song.Song, error)
	GetSongsByScraperID(ctx context.Context, scrID string) ([]song.Song, error)
}
