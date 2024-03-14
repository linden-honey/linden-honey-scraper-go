package scraper

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/linden-honey/linden-honey-api-go/pkg/application/domain/song"
)

// Scraper is an implementation of a songs scraper from some source.
type Scraper struct {
	fetcher Fetcher
	parser  Parser

	songResourcePath             string
	songMetadataListResourcePath string

	validation bool
}

// Fetcher is a component for fetching content in an eager manner.
type Fetcher interface {
	Fetch(ctx context.Context, path string) (string, error)
}

// Parser is a component for parsing content into domain types.
type Parser interface {
	ParseSong(in string) (*song.Entity, error)
	ParseSongMetadataList(in string) ([]song.Metadata, error)
}

// New returns a pointer to the new instance of [Scraper] or an error.
func New(
	f Fetcher,
	p Parser,
	opts ...Option,
) (*Scraper, error) {
	scr := &Scraper{
		fetcher:                      f,
		parser:                       p,
		songResourcePath:             "/%s",
		songMetadataListResourcePath: "/",
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

// WithSongResourcePath sets a templated path to the song resource accepting the song ID as a string for the [Scraper].
func WithSongResourcePath(path string) Option {
	return func(scr *Scraper) {
		scr.songResourcePath = path
	}
}

// WithSongMetadataListResourcePath sets a path to the song metadata list resource for the [Scraper].
func WithSongMetadataListResourcePath(path string) Option {
	return func(scr *Scraper) {
		scr.songMetadataListResourcePath = path
	}
}

// WithValidation enables or disables domain types validation for the [Scraper].
func WithValidation(validation bool) Option {
	return func(scr *Scraper) {
		scr.validation = validation
	}
}

// GetSongs scrapes all songs and returns a slice of [song.Entity] instances or an error.
func (scr *Scraper) GetSongs(ctx context.Context) ([]song.Entity, error) {
	ms, err := scr.GetSongMetadataList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get previews: %w", err)
	}

	outCh := make(chan song.Entity, len(ms))
	errCh := make(chan error, 1)
	for _, m := range ms {
		go func(id string) {
			s, err := scr.GetSong(ctx, id)
			if err != nil {
				errCh <- fmt.Errorf("failed to get a song with id=%s: %w", id, err)
				return
			}

			outCh <- *s
		}(m.ID)
	}

	out := make([]song.Entity, 0)
loop:
	for {
		select {
		case s := <-outCh:
			out = append(out, s)
			if len(out) == len(ms) {
				break loop
			}
		case err := <-errCh:
			return nil, err
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Title < out[j].Title
	})

	return out, nil
}

// GetSongMetadataList scrapes a list of songs and returns a slice of [song.Metadata] instances or an error.
func (scr *Scraper) GetSongMetadataList(ctx context.Context) ([]song.Metadata, error) {
	data, err := scr.fetcher.Fetch(ctx, scr.songMetadataListResourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	out, err := scr.parser.ParseSongMetadataList(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse previews: %w", err)
	}

	errs := make([]error, 0)
	for _, p := range out {
		if scr.validation {
			if err := p.Validate(); err != nil {
				errs = append(errs, fmt.Errorf("failed to validate a preview with id=%s : %w", p.ID, err))
			}
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to validate previews: %w", errors.Join(errs...))
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Title < out[j].Title
	})

	return out, nil
}

// GetSong scrapes a song by id and returns a pointer to the new instance of [song.Entity] or an error.
func (scr *Scraper) GetSong(ctx context.Context, id string) (*song.Entity, error) {
	content, err := scr.fetcher.Fetch(ctx, fmt.Sprintf(scr.songResourcePath, id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	out, err := scr.parser.ParseSong(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse a song: %w", err)
	}

	out.ID = id // backfill ID

	if scr.validation {
		if err := out.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate a song: %w", err)
		}
	}

	return out, nil
}
