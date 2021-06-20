package scraper

import (
	"context"

	"github.com/linden-honey/linden-honey-go/pkg/song"
)

// Service represents the song service interface
type Service interface {
	GetSong(ctx context.Context, id string) (*song.Song, error)
	GetSongs(ctx context.Context) ([]song.Song, error)
	GetPreviews(ctx context.Context) ([]song.Preview, error)
}
