package scraper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gojektech/heimdall"

	"github.com/gojektech/heimdall/httpclient"

	"github.com/linden-honey/linden-honey-scraper-go/internal/util/parser"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/domain"
)

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"

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

func (scraper *scraper) FetchSong(ID string) *domain.Song {
	return nil
}

func (scraper *scraper) FetchSongs() (songs []domain.Song) {
	return nil
}

func (scraper *scraper) FetchPreviews() []domain.Preview {
	previews := make([]domain.Preview, 0)
	url := fmt.Sprintf("%s/texts", scraper.baseURL)
	header := http.Header{
		"User-Agent": []string{userAgent},
	}
	res, err := scraper.client.Get(url, header)
	if err != nil {
		log.Println(err)
		return previews
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return previews
	}
	return parser.ParsePreviews(string(body)) // TODO use reader in parser
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
