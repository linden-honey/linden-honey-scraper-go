package scraper

import (
	"context"
)

// Service represents the scraper service interface
type Service interface {
	GetSong(ctx context.Context, ID string) (*Song, error)
	GetSongs(ctx context.Context) ([]Song, error)
	GetPreviews(ctx context.Context) ([]Preview, error)
}
