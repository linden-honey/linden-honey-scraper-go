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
	ParseSong(input string) (*song.Song, error)
	ParsePreviews(input string) ([]song.Preview, error)
}

// Validator represents the validator interface
type Validator interface {
	Validate(s interface{}) bool
}

// scraper represents the scraper implementation
type scraper struct {
	fetcher   Fetcher
	parser    Parser
	validator Validator
}

// NewService returns a pointer to the new instance of scraper
func NewService(
	fetcher Fetcher,
	parser Parser,
	validator Validator,
) song.Service {
	return &scraper{
		fetcher:   fetcher,
		parser:    parser,
		validator: validator,
	}
}

func (scr *scraper) GetSong(ctx context.Context, id string) (*song.Song, error) {
	data, err := scr.fetcher.Fetch("text_print.php?area=go_texts&id=%s", id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	s, err := scr.parser.ParseSong(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse song: %w", err)
	}

	if !scr.validator.Validate(s) {
		return nil, fmt.Errorf("song %v is invalid", s)
	}

	return s, nil
}

func (scr *scraper) GetSongs(ctx context.Context) ([]song.Song, error) {
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

func (scr *scraper) GetPreviews(_ context.Context) ([]song.Preview, error) {
	data, err := scr.fetcher.Fetch("texts")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	pp, err := scr.parser.ParsePreviews(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse previews: %w", err)
	}

	validPP := make([]song.Preview, 0)
	for _, p := range pp {
		if scr.validator.Validate(p) {
			validPP = append(validPP, p)
		}
	}

	sort.SliceStable(pp, func(i, j int) bool {
		return pp[i].Title < pp[j].Title
	})

	return validPP, nil
}
