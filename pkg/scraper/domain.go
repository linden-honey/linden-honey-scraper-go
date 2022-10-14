package scraper

import (
	"context"

	"github.com/linden-honey/linden-honey-api-go/pkg/song"
)

// Service represents the songs scraper interface
type Service interface {
	GetSong(ctx context.Context, id string) (*song.Song, error)
	GetSongs(ctx context.Context) ([]song.Song, error)
	GetPreviews(ctx context.Context) ([]song.Metadata, error)
}
