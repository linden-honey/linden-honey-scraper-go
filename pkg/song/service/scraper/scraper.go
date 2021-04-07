package scraper

import (
	"context"
	"fmt"
	"sort"

	"github.com/linden-honey/linden-honey-go/pkg/song"
	"github.com/linden-honey/linden-honey-sdk-go/validation"
)

// Fetcher represents the content fetcher interface
type Fetcher interface {
	Fetch(pathFormat string, args ...interface{}) (string, error)
}

// Parser represents the parser interface
type Parser interface {
	ParseSong(in string) (*song.Song, error)
	ParsePreviews(in string) ([]song.Preview, error)
}

// Validator represents the validator interface
type Validator interface {
	Validate(in validation.Validator) error
}

// Scraper represents the default scraper implementation
type Scraper struct {
	fetcher   Fetcher
	parser    Parser
	validator Validator
}

// NewScraper returns a pointer to the new instance of Scraper or an error
func NewScraper(
	fetcher Fetcher,
	parser Parser,
	validator Validator,
) (*Scraper, error) {
	return &Scraper{
		fetcher:   fetcher,
		parser:    parser,
		validator: validator,
	}, nil
}

// GetSong scrapes a song from some source and returns it or an error
func (scr *Scraper) GetSong(_ context.Context, id string) (*song.Song, error) {
	data, err := scr.fetcher.Fetch("text_print.php?area=go_texts&id=%s", id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	s, err := scr.parser.ParseSong(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse a song: %w", err)
	}

	if err := scr.validator.Validate(s); err != nil {
		return nil, fmt.Errorf("failed to validate a song %v: %w", s, err)
	}

	return s, nil
}

// GetSongs scrapes songs from some source and returns them or an error
func (scr *Scraper) GetSongs(ctx context.Context) ([]song.Song, error) {
	previews, err := scr.GetPreviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get previews: %w", err)
	}

	songc := make(chan song.Song, len(previews))
	errc := make(chan error, 1)
	for _, p := range previews {
		go func(id string) {
			s, err := scr.GetSong(ctx, id)
			if err != nil {
				errc <- fmt.Errorf("failed to get a song with id %s: %w", id, err)
				return
			}

			songc <- *s
		}(p.ID)
	}

	songs := make([]song.Song, 0, len(previews))
loop:
	for {
		select {
		case s := <-songc:
			songs = append(songs, s)
			if len(songs) == len(previews) {
				break loop
			}
		case err := <-errc:
			return nil, err
		}
	}

	sort.SliceStable(songs, func(i, j int) bool {
		return songs[i].Title < songs[j].Title
	})

	return songs, nil
}

// GetPreviews scrapes previews from some source and returns them or an error
func (scr *Scraper) GetPreviews(_ context.Context) ([]song.Preview, error) {
	data, err := scr.fetcher.Fetch("texts")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	previews, err := scr.parser.ParsePreviews(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse previews: %w", err)
	}

	for _, p := range previews {
		if err := scr.validator.Validate(p); err != nil {
			return nil, fmt.Errorf("failed to validate a preview %v: %w", p, err)
		}
	}

	sort.SliceStable(previews, func(i, j int) bool {
		return previews[i].Title < previews[j].Title
	})

	return previews, nil
}
