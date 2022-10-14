package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/text/encoding/charmap"
)

// Fetcher represents an implementation of the fetcher
type Fetcher struct {
	baseURL  *url.URL
	encoding *charmap.Charmap
	client   httpClient
	retry    *RetryConfig
}

// RetryConfig represents the retry configuration of the fetcher
type RetryConfig struct {
	Attempts   int
	MinTimeout time.Duration
	MaxTimeout time.Duration
	Factor     time.Duration
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Option set optional parameters for the fetcher
type Option func(*Fetcher)

// NewFetcher returns a pointer to a new instance of the fetcher or an error
func NewFetcher(
	baseURL *url.URL,
	encoding *charmap.Charmap,
	opts ...Option,
) (*Fetcher, error) {
	f := &Fetcher{
		baseURL:  baseURL,
		encoding: encoding,
		client:   new(http.Client),
	}

	for _, opt := range opts {
		opt(f)
	}

	if err := f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

// WithClient sets the http client for the fetcher
func WithClient(client httpClient) Option {
	return func(f *Fetcher) {
		f.client = client
	}
}

// WithRetry sets the retry configuration of the fetcher
func WithRetry(cfg *RetryConfig) Option {
	return func(f *Fetcher) {
		f.retry = cfg
	}
}

// Fetch send GET request under relative path and returns content as a string
func (f *Fetcher) Fetch(ctx context.Context, path string) (string, error) {
	u, err := f.baseURL.Parse(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse an URL: %w", err)
	}

	if f.retry != nil {
		return f.fetchWithRetry(ctx, u)
	}

	return f.fetch(ctx, u)
}

func (f *Fetcher) fetchWithRetry(ctx context.Context, u *url.URL) (string, error) {
	attempt := 0
	for {
		res, err := f.fetch(ctx, u)
		if err != nil {
			if attempt == f.retry.Attempts-1 {
				return "", fmt.Errorf("failed to fetch after attempts=%d: %w", attempt+1, err)
			}

			delay := f.retry.Factor * time.Duration(attempt)
			if delay < f.retry.MinTimeout {
				delay = f.retry.MinTimeout
			} else if delay > f.retry.MinTimeout {
				delay = f.retry.MinTimeout
			}

			select {
			case <-time.After(delay):
				attempt++
				continue
			case <-ctx.Done():
				return "", fmt.Errorf("failed to retry fetch, attempt=%ds: %w", attempt+1, ctx.Err())
			}
		}

		return res, nil
	}
}

func (f *Fetcher) fetch(ctx context.Context, u *url.URL) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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
