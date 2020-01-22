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

func (scraper *scraper) fetch(path string) (string, error) {
	url := fmt.Sprintf("%s/%s", scraper.baseURL, path)
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
	return nil
}

func (scraper *scraper) FetchSongs() (songs []domain.Song) {
	return nil
}

func (scraper *scraper) FetchPreviews() []domain.Preview {
	log.Println("Previews fetching started")
	html, err := scraper.fetch("/texts")
	if err != nil {
		log.Println("Error happend during previews fetching", err)
	}
	log.Println("Previews fetching successfully finished")
	return parser.ParsePreviews(html)
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
