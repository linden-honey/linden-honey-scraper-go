package fetcher

import (
	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
	"time"
)

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
