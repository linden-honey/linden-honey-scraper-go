package scraper

import (
	"context"
	"fmt"
	"sort"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
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
	pp, err := scr.GetPreviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get previews: %w", err)
	}

	sc := make(chan song.Song, len(pp))
	errc := make(chan error, 1)
	for _, p := range pp {
		go func(id string) {
			s, err := scr.GetSong(ctx, id)
			if err != nil {
				errc <- fmt.Errorf("failed to get a song with id %s: %w", id, err)
				return
			}

			sc <- *s
		}(p.ID)
	}

	ss := make([]song.Song, 0, len(pp))
loop:
	for {
		select {
		case s := <-sc:
			ss = append(ss, s)
			if len(ss) == len(pp) {
				break loop
			}
		case err := <-errc:
			return nil, err
		}
	}

	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].Title < ss[j].Title
	})

	return ss, nil
}

// GetPreviews scrapes previews from some source and returns them or an error
func (scr *Scraper) GetPreviews(_ context.Context) ([]song.Preview, error) {
	data, err := scr.fetcher.Fetch("texts")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	pp, err := scr.parser.ParsePreviews(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse previews: %w", err)
	}

	for _, p := range pp {
		if err := scr.validator.Validate(p); err != nil {
			return nil, fmt.Errorf("failed to validate a preview %v: %w", p, err)
		}
	}

	sort.SliceStable(pp, func(i, j int) bool {
		return pp[i].Title < pp[j].Title
	})

	return pp, nil
}
