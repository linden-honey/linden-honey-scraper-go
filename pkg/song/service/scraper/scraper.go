package scraper

import (
	"context"
	"fmt"
	"sort"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/domain"
)

// Fetcher represents the content fetcher interface
type Fetcher interface {
	Fetch(pathFormat string, args ...interface{}) (string, error)
}

// Parser represents the parser interface
type Parser interface {
	ParseSong(input string) (*domain.Song, error)
	ParsePreviews(input string) ([]domain.Preview, error)
}

// Validator represents the validator interface
type Validator interface {
	Validate(s interface{}) bool
}

// Scraper represents the scraper implementation
type Scraper struct {
	fetcher   Fetcher
	parser    Parser
	validator Validator
}

// NewScraper returns a pointer to the new instance of Scraper
func NewScraper(
	fetcher Fetcher,
	parser Parser,
	validator Validator,
) *Scraper {
	return &Scraper{
		fetcher:   fetcher,
		parser:    parser,
		validator: validator,
	}
}

func (s *Scraper) GetSong(ctx context.Context, id string) (*domain.Song, error) {
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

func (s *Scraper) GetSongs(ctx context.Context) ([]domain.Song, error) {
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

func (s *Scraper) GetPreviews(_ context.Context) ([]domain.Preview, error) {
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
