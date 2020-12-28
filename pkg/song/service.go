package song

import (
	"context"
)

// Service represents a song service interface
type Service interface {
	GetSong(ctx context.Context, id string) (*Song, error)
	GetSongs(ctx context.Context) ([]Song, error)
	GetPreviews(ctx context.Context) ([]Preview, error)
}

type Middleware func(Service) Service
