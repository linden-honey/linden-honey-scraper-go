package scraper

import (
	"context"
	"fmt"
	"sort"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
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
	ValidateSong(s song.Song) error
	ValidatePreview(s song.Preview) error
}

// Scraper represents the default scraper implementation
type Scraper struct {
	fetcher   Fetcher
	parser    Parser
	validator Validator
}

// NewScraper returns a pointer to the new instance of scraper
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

func (scr *Scraper) GetSong(_ context.Context, id string) (*song.Song, error) {
	data, err := scr.fetcher.Fetch("text_print.php?area=go_texts&id=%s", id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	s, err := scr.parser.ParseSong(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse song: %w", err)
	}

	if err := scr.validator.ValidateSong(*s); err != nil {
		return nil, fmt.Errorf("song %v is invalid: %w", s, err)
	}

	return s, nil
}

func (scr *Scraper) GetSongs(ctx context.Context) ([]song.Song, error) {
	pp, err := scr.GetPreviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get previews: %w", err)
	}

	ss := make([]song.Song, 0)
	for _, p := range pp {
		s, err := scr.GetSong(ctx, p.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get song with id %s", p.ID)
		}
		ss = append(ss, *s)
	}

	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].Title < ss[j].Title
	})

	return ss, nil
}

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
		if err := scr.validator.ValidatePreview(p); err != nil {
			return nil, fmt.Errorf("preview %v is invalid", p)
		}
	}

	sort.SliceStable(pp, func(i, j int) bool {
		return pp[i].Title < pp[j].Title
	})

	return pp, nil
}
