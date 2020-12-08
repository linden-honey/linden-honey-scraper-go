package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
	"golang.org/x/text/encoding/charmap"
)

//TODO refactor configuration
// 1. rename Properties to Config
// 2. Provide default values
// 3. Validate required values in constructor

// RetryProperties represents the retry properties structure
type RetryProperties struct {
	Retries    int
	Factor     float64
	MinTimeout time.Duration
	MaxTimeout time.Duration
}

// Properties represents the defaultScraper properties structure
type Properties struct {
	BaseURL        *url.URL
	SourceEncoding *charmap.Charmap
}

// defaultFetcher represents the default fetcher implementation
type defaultFetcher struct {
	client *httpclient.Client
	props  *Properties
}

// NewDefaultFetcher returns a pointer to the new instance of defaultFetcher
func NewDefaultFetcher(props *Properties) *defaultFetcher {
	return &defaultFetcher{
		client: httpclient.NewClient(),
		props:  props,
	}
}

// NewDefaultFetcherWithRetry returns pointer to the new instance of defaultFetcher with retry feature
func NewDefaultFetcherWithRetry(props *Properties, retry *RetryProperties) *defaultFetcher {
	return &defaultFetcher{
		client: httpclient.NewClient(
			httpclient.WithRetryCount(retry.Retries),
			httpclient.WithRetrier(
				heimdall.NewRetrier(
					heimdall.NewExponentialBackoff(
						retry.MinTimeout,
						retry.MaxTimeout,
						retry.Factor,
						time.Second,
					),
				),
			),
		),
		props: props,
	}
}

// Fetch send GET request under relative path built with pathFormat and args and returns content string
func (f *defaultFetcher) Fetch(pathFormat string, args ...interface{}) (string, error) {
	fetchURL, err := f.props.BaseURL.Parse(fmt.Sprintf(pathFormat, args...))
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	header := http.Header{
		"User-Agent": []string{
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		},
	}
	res, err := f.client.Get(fetchURL.String(), header)
	if err != nil {
		return "", fmt.Errorf("failed to proceed request: %w", err)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("server did not respond successfully - status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	decoder := f.props.SourceEncoding.NewDecoder()
	body, err := ioutil.ReadAll(decoder.Reader(res.Body))
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}
