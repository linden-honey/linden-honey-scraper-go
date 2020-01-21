package scraper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/linden-honey/linden-honey-scraper-go/internal/util/parser"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/domain"
)

// RetryProperties represents a retry properties structure
type RetryProperties struct {
}

// Properties represents a scraper properties structure
type Properties struct {
	BaseURL         string
	RetryProperties *RetryProperties
}

// Scraper represents a scraper interface
type Scraper interface {
	FetchSong(ID string) *domain.Song
	FetchSongs() []domain.Song
	FetchPreviews() []domain.Preview
}

type scraper struct {
	Properties *Properties
}

func (scraper *scraper) FetchSong(ID string) *domain.Song {
	return nil
}

func (scraper *scraper) FetchSongs() (songs []domain.Song) {
	return nil
}

func (scraper *scraper) FetchPreviews() []domain.Preview {
	previews := make([]domain.Preview, 0)
	url := fmt.Sprintf("%s/texts", scraper.Properties.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return previews
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return previews
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	stringBody := string(body)
	log.Println(stringBody)
	if err != nil {
		log.Println(err)
		return previews
	}
	return parser.ParsePreviews(stringBody)
}

// Create returns a scraper instance
func Create(properties *Properties) Scraper {
	return &scraper{
		Properties: properties,
	}
}
