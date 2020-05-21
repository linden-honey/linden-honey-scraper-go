package scraper

import (
	"sort"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/domain"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/fetcher"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/parser"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/validator"
)

// Scraper represents the scraper interface
type Scraper interface {
	GetSong(ID string) (*domain.Song, error)
	GetSongs() ([]*domain.Song, error)
	GetPreviews() ([]*domain.Preview, error)
}

// defaultScraper represents the default scraper implementation
type defaultScraper struct {
	logger    *log.Logger
	fetcher   fetcher.Fetcher
	parser    parser.Parser
	validator validator.Validator
}

// NewDefaultScraper returns a pointer to the new instance of defaultScraper
func NewDefaultScraper(logger *log.Logger, fetcher fetcher.Fetcher, parser parser.Parser, validator validator.Validator) *defaultScraper {
	return &defaultScraper{
		logger:    logger,
		fetcher:   fetcher,
		parser:    parser,
		validator: validator,
	}
}

// GetSong returns pointer to the scrapped Song instance
func (s *defaultScraper) GetSong(ID string) (*domain.Song, error) {
	s.logger.Infof("Fetching song with id %s", ID)
	html, err := s.fetcher.Fetch("text_print.php?area=go_texts&id=%s", ID)
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
	s.logger.Infof(`Successfully fetched song with id %s and title "%s"`, ID, song.Title)
	return song, nil
}

// GetSongs returns a slice of pointers of the scrapped Song instances
func (s *defaultScraper) GetSongs() ([]*domain.Song, error) {
	s.logger.Info("Songs fetching started")
	previews, err := s.GetPreviews()
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during previews fetching")
	}
	var wg sync.WaitGroup
	var m sync.Mutex
	songs := make([]*domain.Song, 0)
	for _, p := range previews {
		wg.Add(1)
		go func() {
			defer wg.Done()
			song, err := s.GetSong(p.ID)
			if err != nil {
				s.logger.Warn(errors.Wrapf(err, "Error happened during fetching song with id %s"))
			} else {
				m.Lock()
				songs = append(songs, song)
				m.Unlock()
			}
		}()
	}
	wg.Wait()
	sort.SliceStable(songs, func(i, j int) bool {
		return songs[i].Title < songs[j].Title
	})
	s.logger.Info("Songs fetching successfully finished")
	return songs, nil
}

// GetSongs returns a slice of pointers of the scrapped Preview instances
func (s *defaultScraper) GetPreviews() ([]*domain.Preview, error) {
	s.logger.Info("Previews fetching started")
	html, err := s.fetcher.Fetch("texts")
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during previews html fetching")
	}
	previews, err := s.parser.ParsePreviews(html)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during previews html parsing")
	}
	validPreviews := make([]*domain.Preview, 0)
	for _, preview := range previews {
		if s.validator.Validate(preview) {
			validPreviews = append(validPreviews, preview)
		}
	}
	sort.SliceStable(previews, func(i, j int) bool {
		return previews[i].Title < previews[j].Title
	})
	s.logger.Info("Previews fetching successfully finished")
	return validPreviews, nil
}
