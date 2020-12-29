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

	"github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Fetcher represents the default fetcher implementation
type Fetcher struct {
	client         heimdall.Doer
	baseURL        *url.URL
	sourceEncoding *charmap.Charmap
}

// Config represents common the fetcher configuration
type Config struct {
	BaseURL        *url.URL
	SourceEncoding *charmap.Charmap
}

func (cfg *Config) Validate() error {
	if cfg.BaseURL == nil {
		return errors.NewRequiredValueError("BaseURL")
	}

	if cfg.SourceEncoding == nil {
		return errors.NewRequiredValueError("SourceEncoding")
	}

	return nil
}

// RetryConfig represents the retry configuration
type RetryConfig struct {
	Retries           int
	Factor            float64
	MinTimeout        time.Duration
	MaxTimeout        time.Duration
	MaxJitterInterval time.Duration
}

// NewFetcher returns a pointer to the new instance of Fetcher
func NewFetcher(cfg Config) (*Fetcher, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid fetcher config: %w", err)
	}

	return &Fetcher{
		client:         httpclient.NewClient(),
		baseURL:        cfg.BaseURL,
		sourceEncoding: cfg.SourceEncoding,
	}, nil
}

// NewFetcherWithRetry returns a pointer to the new instance of defaultFetcher with retry feature
func NewFetcherWithRetry(cfg Config, retryCfg RetryConfig) (*Fetcher, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid fetcher config: %w", err)
	}

	return &Fetcher{
		client: httpclient.NewClient(
			httpclient.WithRetryCount(retryCfg.Retries),
			httpclient.WithRetrier(
				heimdall.NewRetrier(
					heimdall.NewExponentialBackoff(
						retryCfg.MinTimeout,
						retryCfg.MaxTimeout,
						retryCfg.Factor,
						retryCfg.MaxJitterInterval,
					),
				),
			),
		),
		baseURL:        cfg.BaseURL,
		sourceEncoding: cfg.SourceEncoding,
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
