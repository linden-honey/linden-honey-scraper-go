package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/text/encoding"
)

// Fetcher is an implementation of an eager content fetcher.
type Fetcher struct {
	baseURL  url.URL
	client   httpClient
	encoding encoding.Encoding
	retry    retryFunc
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type retryFunc func(ctx context.Context, action func() (string, error)) (string, error)

// New returns a pointer to the new instance of [Fetcher] or an error.
func New(
	baseURL url.URL,
	opts ...Option,
) (*Fetcher, error) {
	f := &Fetcher{
		baseURL: baseURL,
		client:  new(http.Client),
	}

	for _, opt := range opts {
		opt(f)
	}

	if err := f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

// Option set optional parameters for the [Fetcher].
type Option func(*Fetcher)

// WithClient sets the http client for the [Fetcher].
func WithClient(client httpClient) Option {
	return func(f *Fetcher) {
		f.client = client
	}
}

// WithEncoding sets the encoding for the [Fetcher].
func WithEncoding(encoding encoding.Encoding) Option {
	return func(f *Fetcher) {
		f.encoding = encoding
	}
}

// WithRetry sets the retry function for the [Fetcher].
func WithRetry(retry retryFunc) Option {
	return func(f *Fetcher) {
		f.retry = retry
	}
}

// Fetch gets a resource under the relative path and returns it as a UTF-8 encoded string or returns an error.
func (f *Fetcher) Fetch(ctx context.Context, path string) (string, error) {
	u, err := f.baseURL.Parse(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse an URL: %w", err)
	}

	if f.retry != nil {
		return f.retry(ctx, func() (string, error) {
			return f.fetch(ctx, u)
		})
	}

	return f.fetch(ctx, u)
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

	defer res.Body.Close()

	var r io.Reader = res.Body
	if f.encoding != nil {
		r = f.encoding.NewDecoder().Reader(r)
	}

	body, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to read a response: %w", err)
	}

	return string(body), nil
}
