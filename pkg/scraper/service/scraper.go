package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/domain"
)

type fetcher interface {
	Fetch(pathFormat string, args ...interface{}) (string, error)
}

type parser interface {
	ParseSong(input string) (*domain.Song, error)
	ParsePreviews(input string) ([]domain.Preview, error)
}

type validator interface {
	Validate(s interface{}) bool
}

type scraper struct {
	fetcher   fetcher
	parser    parser
	validator validator
}

func NewScraper(
	fetcher fetcher,
	parser parser,
	validator validator,
) *scraper {
	return &scraper{
		fetcher:   fetcher,
		parser:    parser,
		validator: validator,
	}
}

func (s *scraper) GetSong(ctx context.Context, id string) (*domain.Song, error) {
	data, err := s.fetcher.Fetch("text_print.php?area=go_texts&id=%s", id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	song, err := s.parser.ParseSong(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse song: %w", err)
	}

	if !s.validator.Validate(song) {
		return nil, fmt.Errorf("song %v is invalid", song)
	}

	return song, nil
}

func (s *scraper) GetSongs(ctx context.Context) ([]domain.Song, error) {
	previews, err := s.GetPreviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get previews: %w", err)
	}

	songs := make([]domain.Song, 0)
	for _, p := range previews {
		song, err := s.GetSong(ctx, p.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get song with id %s", p.ID)
		}
		songs = append(songs, *song)
	}

	sort.SliceStable(songs, func(i, j int) bool {
		return songs[i].Title < songs[j].Title
	})

	return songs, nil
}

func (s *scraper) GetPreviews(_ context.Context) ([]domain.Preview, error) {
	data, err := s.fetcher.Fetch("texts")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	previews, err := s.parser.ParsePreviews(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse previews: %w", err)
	}

	validPreviews := make([]domain.Preview, 0)
	for _, p := range previews {
		if s.validator.Validate(p) {
			validPreviews = append(validPreviews, p)
		}
	}

	sort.SliceStable(previews, func(i, j int) bool {
		return previews[i].Title < previews[j].Title
	})

	return validPreviews, nil
}
