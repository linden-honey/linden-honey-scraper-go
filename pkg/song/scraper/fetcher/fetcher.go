package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
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

// Properties represents the properties structure
type Properties struct {
	BaseURL        *url.URL
	SourceEncoding *charmap.Charmap
}

// Fetcher represents the default fetcher implementation
type Fetcher struct {
	client         heimdall.Doer
	baseURL        *url.URL
	sourceEncoding *charmap.Charmap
}

// NewFetcher returns a pointer to the new instance of Fetcher
func NewFetcher(props *Properties) (*Fetcher, error) {
	return &Fetcher{
		client:         httpclient.NewClient(),
		baseURL:        props.BaseURL,
		sourceEncoding: props.SourceEncoding,
	}, nil
}

// NewFetcherWithRetry returns a pointer to the new instance of defaultFetcher with retry feature
func NewFetcherWithRetry(props *Properties, retry *RetryProperties) (*Fetcher, error) {
	return &Fetcher{
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
		baseURL:        props.BaseURL,
		sourceEncoding: props.SourceEncoding,
	}, nil
}

// Fetch send GET request under relative path built with pathFormat and args and returns content string
func (f *Fetcher) Fetch(pathFormat string, args ...interface{}) (string, error) {
	fetchURL, err := f.baseURL.Parse(fmt.Sprintf(pathFormat, args...))
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, fetchURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = http.Header{
		"User-Agent": []string{
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		},
	}

	res, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to proceed request: %w", err)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("server did not respond successfully - status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	decoder := f.sourceEncoding.NewDecoder()
	body, err := ioutil.ReadAll(decoder.Reader(res.Body))
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}
