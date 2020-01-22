package scraper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"

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
	BaseURL         string
	RetryProperties RetryProperties
}

// Scraper represents a scraper interface
type Scraper interface {
	FetchSong(ID string) *domain.Song
	FetchSongs() []domain.Song
	FetchPreviews() []domain.Preview
}

type scraper struct {
	baseURL string
	client  *httpclient.Client
}

func (scraper *scraper) fetch(path string, args ...interface{}) (string, error) {
	url := fmt.Sprintf("%s/%s", scraper.baseURL, fmt.Sprintf(path, args...))
	header := http.Header{
		"User-Agent": []string{
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		},
	}
	res, err := scraper.client.Get(url, header)
	if err != nil {
		return "", errors.Wrap(err, "Fetch operation failed")
	}
	if res.StatusCode != 200 {
		return "", errors.Wrapf(err, "Server did not respond successfully - status code %d", res.StatusCode)
	}
	defer res.Body.Close()
	decoder := charmap.Windows1251.NewDecoder()
	body, err := ioutil.ReadAll(decoder.Reader(res.Body))
	if err != nil {
		return "", errors.Wrap(err, "Error during responce reading")
	}
	return string(body), nil
}

func (scraper *scraper) FetchSong(ID string) *domain.Song {
	log.Printf("Fetching song with id %s", ID)
	html, err := scraper.fetch("text_print.php?area=go_texts&id=%s", ID)
	if err != nil {
		log.Printf("Error happend during fetching song with id %s", ID)
		log.Println(err)
	}
	song := parser.ParseSong(html)
	if !validation.Validate(song) {
		return nil
	}
	log.Printf(`Successfully fetched song with id %s and title "%s"`, ID, song.Title)
	return song
}

func (scraper *scraper) FetchSongs() []domain.Song {
	log.Println("Songs fetching started")
	previews := scraper.FetchPreviews()
	songs := make([]domain.Song, 0)
	for _, preview := range previews {
		song := scraper.FetchSong(preview.ID)
		if song != nil {
			songs = append(songs, *song)
		}
	}
	log.Println("Songs fetching successfully finished")
	return songs
}

func (scraper *scraper) FetchPreviews() []domain.Preview {
	log.Println("Previews fetching started")
	html, err := scraper.fetch("texts")
	previews := make([]domain.Preview, 0)
	if err != nil {
		log.Println("Error happend during previews fetching", err)
		return previews
	}
	for _, preview := range parser.ParsePreviews(html) {
		if validation.Validate(preview) {
			previews = append(previews, preview)
		}
	}
	log.Println("Previews fetching successfully finished")
	return previews
}

// Create returns a scraper instance
func Create(properties *Properties) Scraper {
	return &scraper{
		baseURL: properties.BaseURL,
		client: httpclient.NewClient(
			httpclient.WithRetryCount(properties.RetryProperties.Retries),
			httpclient.WithRetrier(
				heimdall.NewRetrier(
					heimdall.NewExponentialBackoff(
						properties.RetryProperties.MinTimeout,
						properties.RetryProperties.MaxTimeout,
						properties.RetryProperties.Factor,
						time.Second,
					),
				),
			),
		),
	}
}
