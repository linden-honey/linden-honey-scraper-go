package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"golang.org/x/text/encoding/charmap"

	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

// Fetcher represents the default fetcher implementation
type Fetcher struct {
	client         heimdall.Doer
	baseURL        *url.URL
	sourceEncoding *charmap.Charmap
}

// Config represents common Fetcher configuration
type Config struct {
	BaseURL        *url.URL
	SourceEncoding *charmap.Charmap
}

// Validate validates Fetcher configuration
func (cfg *Config) Validate() error {
	if cfg.BaseURL == nil {
		return sdkerrors.NewRequiredValueError("BaseURL")
	}

	if cfg.SourceEncoding == nil {
		return sdkerrors.NewRequiredValueError("SourceEncoding")
	}

	return nil
}

// NewFetcher returns a pointer to the new instance of Fetcher or an error
func NewFetcher(cfg Config, opts ...Option) (*Fetcher, error) {
	if err := cfg.Validate(); err != nil {
		return nil, sdkerrors.NewInvalidValueError("Config", err)
	}

	f := &Fetcher{
		client:         httpclient.NewClient(),
		baseURL:        cfg.BaseURL,
		sourceEncoding: cfg.SourceEncoding,
	}

	for _, opt := range opts {
		opt(f)
	}

	return f, nil
}

type Option func(*Fetcher)

// RetryConfig represents the retry configuration
type RetryConfig struct {
	Retries           int
	Factor            float64
	MinTimeout        time.Duration
	MaxTimeout        time.Duration
	MaxJitterInterval time.Duration
}

func WithRetry(cfg RetryConfig) Option {
	return func(f *Fetcher) {
		// TODO: rewrite with simple attempts count and time.Sleep retry
		f.client = httpclient.NewClient(
			httpclient.WithRetryCount(cfg.Retries),
			httpclient.WithRetrier(
				heimdall.NewRetrier(
					heimdall.NewExponentialBackoff(
						cfg.MinTimeout,
						cfg.MaxTimeout,
						cfg.Factor,
						cfg.MaxJitterInterval,
					),
				),
			),
		)
	}
}

// Fetch send GET request under relative path and returns content as a string
func (f *Fetcher) Fetch(path string) (string, error) {
	fetchURL, err := f.baseURL.Parse(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse an URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, fetchURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create a request: %w", err)
	}

	res, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to proceed request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server did not respond successfully - status code %d", res.StatusCode)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	decoder := f.sourceEncoding.NewDecoder()
	body, err := ioutil.ReadAll(decoder.Reader(res.Body))
	if err != nil {
		return "", fmt.Errorf("failed to read a response: %w", err)
	}

	return string(body), nil
}
