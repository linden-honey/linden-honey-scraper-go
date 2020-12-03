package scraper

import (
	"context"
	"sort"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

// Fetcher represents the fetcher interface
type Fetcher interface {
	Fetch(pathFormat string, args ...interface{}) (string, error)
}

// Parser represents the fetcher interface
type Parser interface {
	ParseSong(html string) (*Song, error)
	ParsePreviews(html string) ([]Preview, error)
}

// Validator represents the validator interface
type Validator interface {
	Validate(s interface{}) bool
}

type scraperService struct {
	logger log.Logger

	fetcher   Fetcher
	parser    Parser
	validator Validator
}

func NewScraperService(
	logger log.Logger,
	fetcher Fetcher,
	parser Parser,
	validator Validator,
) Service {
	return &scraperService{
		logger: logger,

		fetcher:   fetcher,
		parser:    parser,
		validator: validator,
	}
}

func (s scraperService) GetSong(ctx context.Context, id string) (*Song, error) {
	html, err := s.fetcher.Fetch("text_print.php?area=go_texts&id=%s", id)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during song html fetching")
	}
	song, err := s.parser.ParseSong(html)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during song html parsing")
	}
	if !s.validator.Validate(song) {
		return nil, errors.Errorf("Song %V is invalid", song)
	}
	return song, nil
}

func (s scraperService) GetSongs(ctx context.Context) ([]Song, error) {
	previews, err := s.GetPreviews(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during previews fetching")
	}
	var wg sync.WaitGroup
	var m sync.Mutex
	songs := make([]Song, 0)
	for _, p := range previews {
		wg.Add(1)
		go func(p Preview) {
			defer wg.Done()
			song, err := s.GetSong(ctx, p.ID)
			if err != nil {
				_ = level.Warn(s.logger).Log(
					"error", errors.Wrapf(err, "Error happened during fetching song with id %s", p.ID),
				)
			} else {
				m.Lock()
				songs = append(songs, *song)
				m.Unlock()
			}
		}(p)
	}
	wg.Wait()
	sort.SliceStable(songs, func(i, j int) bool {
		return songs[i].Title < songs[j].Title
	})
	return songs, nil
}

func (s scraperService) GetPreviews(ctx context.Context) ([]Preview, error) {
	html, err := s.fetcher.Fetch("texts")
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during previews html fetching")
	}
	previews, err := s.parser.ParsePreviews(html)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during previews html parsing")
	}
	validPreviews := make([]Preview, 0)
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
