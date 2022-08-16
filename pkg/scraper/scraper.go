package scraper

import (
	"context"
	"fmt"
	"sort"

	"github.com/linden-honey/linden-honey-go/pkg/song"
	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Fetcher represents the content fetcher interface
type Fetcher interface {
	Fetch(ctx context.Context, path string) (string, error)
}

// Parser represents the parser interface
type Parser interface {
	ParseSong(in string) (*song.Song, error)
	ParsePreviews(in string) ([]song.Meta, error)
}

// Scraper represents default scraper implementation
type Scraper struct {
	f Fetcher
	p Parser

	validation bool
}

// Option set optional parameters for the scraper
type Option func(*Scraper)

// NewScraper returns a pointer to the new instance of scraper or an error
func NewScraper(
	f Fetcher,
	p Parser,
	opts ...Option,
) (*Scraper, error) {
	scr := &Scraper{
		f: f,
		p: p,
	}

	for _, opt := range opts {
		opt(scr)
	}

	if err := scr.Validate(); err != nil {
		return nil, err
	}

	return scr, nil
}

// WithValidation enables or disables validation
func WithValidation(validation bool) Option {
	return func(scr *Scraper) {
		scr.validation = validation
	}
}

// Validate validates scraper configuration
func (scr *Scraper) Validate() error {
	if scr.f == nil {
		return sdkerrors.NewRequiredValueError("f")
	}

	if scr.p == nil {
		return sdkerrors.NewRequiredValueError("p")
	}

	return nil
}

// GetSong scrapes a song from some source and returns it or an error
func (scr *Scraper) GetSong(ctx context.Context, id string) (*song.Song, error) {
	data, err := scr.f.Fetch(ctx, fmt.Sprintf("text_print.php?area=go_texts&id=%s", id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	s, err := scr.p.ParseSong(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse a song: %w", err)
	}

	if scr.validation {
		if err := s.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate a song %v: %w", s, err)
		}
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
func (scr *Scraper) GetPreviews(ctx context.Context) ([]song.Meta, error) {
	data, err := scr.f.Fetch(ctx, "texts")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	previews, err := scr.p.ParsePreviews(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse previews: %w", err)
	}

	for _, p := range previews {
		if scr.validation {
			if err := p.Validate(); err != nil {
				return nil, fmt.Errorf("failed to validate a preview %v: %w", p, err)
			}
		}
	}

	sort.SliceStable(previews, func(i, j int) bool {
		return previews[i].Title < previews[j].Title
	})

	return previews, nil
}
