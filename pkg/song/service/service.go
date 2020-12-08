package service

import (
	"context"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/domain"
)

// Service represents a scraper interface
type Service interface {
	GetSong(ctx context.Context, id string) (*domain.Song, error)
	GetSongs(ctx context.Context) ([]domain.Song, error)
	GetPreviews(ctx context.Context) ([]domain.Preview, error)
}
