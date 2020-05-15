package scraper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
	"github.com/pkg/errors"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/domain"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/util/parser"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/util/validation"
)

// RetryProperties represents a retry properties structure
type RetryProperties struct {
	Retries    int
	Factor     float64
	MinTimeout time.Duration
	MaxTimeout time.Duration
}

// Properties represents a scraper properties structure
type Properties struct {
	BaseURL *url.URL
	Retry   RetryProperties
}

// Scraper represents a scraper interface
type Scraper interface {
	FetchSong(ID string) (*domain.Song, error)
	FetchSongs() ([]*domain.Song, error)
	FetchPreviews() ([]*domain.Preview, error)
}

type scraper struct {
	baseURL *url.URL
	client  *httpclient.Client
}

func (scraper *scraper) fetch(path string, args ...interface{}) (string, error) {
	fetchURL, err := scraper.baseURL.Parse(fmt.Sprintf(path, args...))
	if err != nil {
		return "", errors.Wrap(err, "Couldn't parse url")
	}
	header := http.Header{
		"User-Agent": []string{
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		},
	}
	res, err := scraper.client.Get(fetchURL.String(), header)
	if err != nil {
		return "", errors.Wrap(err, "GET request failed")
	}
	if res.StatusCode != 200 {
		return "", errors.Wrapf(err, "Server did not respond successfully - status code %d", res.StatusCode)
	}
	defer res.Body.Close()
	decoder := charmap.Windows1251.NewDecoder()
	body, err := ioutil.ReadAll(decoder.Reader(res.Body))
	if err != nil {
		return "", errors.Wrap(err, "Error happened during response reading")
	}
	return string(body), nil
}

func (scraper *scraper) FetchSong(ID string) (*domain.Song, error) {
	log.Printf("Fetching song with id %s", ID)
	html, err := scraper.fetch("text_print.php?area=go_texts&id=%s", ID)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during song html fetching")
	}
	song, err := parser.ParseSong(html)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during song html parsing")
	}
	if !validation.Validate(song) {
		return nil, errors.Errorf("Song %V is invalid", song)
	}
	log.Printf(`Successfully fetched song with id %s and title "%s"`, ID, song.Title)
	return song, nil
}

func (scraper *scraper) FetchSongs() ([]*domain.Song, error) {
	log.Println("Songs fetching started")
	previews, err := scraper.FetchPreviews()
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
			song, err := scraper.FetchSong(p.ID)
			if err != nil {
				log.Println(errors.Wrap(err, "Error happened during song fetching"))
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
	log.Println("Songs fetching successfully finished")
	return songs, nil
}

func (scraper *scraper) FetchPreviews() ([]*domain.Preview, error) {
	log.Println("Previews fetching started")
	html, err := scraper.fetch("texts")
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during previews html fetching")
	}
	previews, err := parser.ParsePreviews(html)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during previews html parsing")
	}
	validPreviews := make([]*domain.Preview, 0)
	for _, preview := range previews {
		if validation.Validate(preview) {
			validPreviews = append(validPreviews, preview)
		} else {
			log.Println("Preview %V is invalid", preview)
		}
	}
	sort.SliceStable(previews, func(i, j int) bool {
		return previews[i].Title < previews[j].Title
	})
	log.Println("Previews fetching successfully finished")
	return validPreviews, nil
}

// NewScraper returns a new scraper instance
func NewScraper(properties *Properties) Scraper {
	return &scraper{
		baseURL: properties.BaseURL,
		client: httpclient.NewClient(
			httpclient.WithRetryCount(properties.Retry.Retries),
			httpclient.WithRetrier(
				heimdall.NewRetrier(
					heimdall.NewExponentialBackoff(
						properties.Retry.MinTimeout,
						properties.Retry.MaxTimeout,
						properties.Retry.Factor,
						time.Second,
					),
				),
			),
		),
	}
}
