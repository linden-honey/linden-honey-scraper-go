package scraper

import (
	"context"
	"fmt"
	"sort"

	"github.com/linden-honey/linden-honey-api-go/pkg/song"
)

// Scraper is an implementation of a song scraper from some source.
type Scraper struct {
	fetcher    Fetcher
	parser     Parser
	validation bool
}

// Fetcher is a component for fetching content in an eager manner.
type Fetcher interface {
	Fetch(ctx context.Context, path string) (string, error)
}

// Parser is a component for parsing content into domain types.
type Parser interface {
	ParseSong(input string) (*song.Song, error)
	ParsePreviews(input string) ([]song.Metadata, error)
}

// New returns a pointer to the new instance of [Scraper] or an error.
func New(
	f Fetcher,
	p Parser,
	opts ...Option,
) (*Scraper, error) {
	scr := &Scraper{
		fetcher: f,
		parser:  p,
	}

	for _, opt := range opts {
		opt(scr)
	}

	if err := scr.Validate(); err != nil {
		return nil, err
	}

	return scr, nil
}

// Option set optional parameters for the [Scraper].
type Option func(*Scraper)

// WithValidation enables or disables domain types validation for the [Scraper].
func WithValidation(validation bool) Option {
	return func(scr *Scraper) {
		scr.validation = validation
	}
}

// GetSong scrapes a song by id and returns a pointer to the new instance of [song.Song] or an error.
func (scr *Scraper) GetSong(ctx context.Context, id string) (*song.Song, error) {
	data, err := scr.fetcher.Fetch(ctx, fmt.Sprintf("text_print.php?area=go_texts&id=%s", id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	s, err := scr.parser.ParseSong(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse a song: %w", err)
	}

	s.ID = id // backfill ID

	if scr.validation {
		if err := s.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate a song: %w", err)
		}
	}

	return s, nil
}

// GetSongs scrapes all songs and returns a slice of [song.Song] instances or an error.
func (scr *Scraper) GetSongs(ctx context.Context) ([]song.Song, error) {
	ps, err := scr.GetPreviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get previews: %w", err)
	}

	sc := make(chan song.Song, len(ps))
	errc := make(chan error, 1)
	for _, p := range ps {
		go func(id string) {
			s, err := scr.GetSong(ctx, id)
			if err != nil {
				errc <- fmt.Errorf("failed to get a song with id=%s: %w", id, err)
				return
			}

			sc <- *s
		}(p.ID)
	}

	ss := make([]song.Song, 0, len(ps))
loop:
	for {
		select {
		case s := <-sc:
			ss = append(ss, s)
			if len(ss) == len(ps) {
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

// GetPreviews scrapes songs metadata and returns a slice of [song.Metadata] instances or an error.
func (scr *Scraper) GetPreviews(ctx context.Context) ([]song.Metadata, error) {
	data, err := scr.fetcher.Fetch(ctx, "texts")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	ps, err := scr.parser.ParsePreviews(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse previews: %w", err)
	}

	for _, p := range ps {
		if scr.validation {
			if err := p.Validate(); err != nil {
				return nil, fmt.Errorf("failed to validate a preview with id=%s : %w", p.ID, err)
			}
		}
	}

	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].Title < ps[j].Title
	})

	return ps, nil
}
