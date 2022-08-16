package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
	"golang.org/x/text/encoding/charmap"

	sdkerrors "github.com/linden-honey/linden-honey-sdk-go/errors"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Fetcher represents default fetcher implementation
type Fetcher struct {
	client httpClient

	baseURL  *url.URL
	encoding *charmap.Charmap
}

// Option set optional parameters for the fetcher
type Option func(*Fetcher)

// NewFetcher returns a pointer to the new instance of fetcher or an error
func NewFetcher(
	baseURL *url.URL,
	encoding *charmap.Charmap,
	opts ...Option,
) (*Fetcher, error) {
	f := &Fetcher{
		client:   httpclient.NewClient(),
		baseURL:  baseURL,
		encoding: encoding,
	}

	for _, opt := range opts {
		opt(f)
	}

	if err := f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

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

// Validate validates fetcher configuration
func (f *Fetcher) Validate() error {
	if f.client == nil {
		return sdkerrors.NewRequiredValueError("client")
	}

	if f.baseURL == nil {
		return sdkerrors.NewRequiredValueError("baseURL")
	}

	if f.encoding == nil {
		return sdkerrors.NewRequiredValueError("encoding")
	}

	return nil
}

// Fetch send GET request under relative path and returns content as a string
func (f *Fetcher) Fetch(ctx context.Context, path string) (string, error) {
	fetchURL, err := f.baseURL.Parse(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse an URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fetchURL.String(), nil)
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

	decoder := f.encoding.NewDecoder()
	body, err := io.ReadAll(decoder.Reader(res.Body))
	if err != nil {
		return "", fmt.Errorf("failed to read a response: %w", err)
	}

	return string(body), nil
}
