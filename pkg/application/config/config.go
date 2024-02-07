package config

import (
	"fmt"
	"net/url"
	"reflect"
	"time"

	"golang.org/x/text/encoding"

	"github.com/caarlos0/env/v6"
)

// Config is a configuration object.
type Config struct {
	Scrapers ScrapersConfig
	Output   OutputConfig
}

// ScrapersConfig is a configuration object.
type ScrapersConfig struct {
	Grob ScraperConfig `envPrefix:"GROB_"`
}

// ScraperConfig is a configuration object.
type ScraperConfig struct {
	BaseURL    url.URL           `env:"SCRAPER_BASE_URL"`
	Encoding   encoding.Encoding `env:"SCRAPER_ENCODING"`
	Validation bool              `env:"SCRAPER_VALIDATION"`
	Retry      RetryConfig       `envPrefix:"SCRAPER_"`
}

// OutputConfig is a configuration object.
type OutputConfig struct {
	FileName string `env:"OUTPUT_FILE_NAME"`
}

// RetryConfig is a configuration object.
type RetryConfig struct {
	Enabled        bool          `env:"RETRY_ENABLED"`
	Attempts       uint          `env:"RETRY_ATTEMPTS"`
	MinInterval    time.Duration `env:"RETRY_MIN_INTERVAL"`
	MaxInterval    time.Duration `env:"RETRY_MAX_INTERVAL"`
	Factor         float64       `env:"RETRY_FACTOR"`
	MaxElapsedTime time.Duration `env:"RETRY_MAX_ELAPSED_TIME"`
}

// New returns a pointer to the new instance of [Config] or an error.
func New() (*Config, error) {
	cfg := Default()

	fmt.Println(reflect.TypeOf((*encoding.Encoding)(nil)).Elem())

	if err := env.ParseWithFuncs(&cfg, Parsers(), env.Options{
		Prefix: "APPLICATION_",
	}); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate a config: %w", err)
	}

	return &cfg, nil
}
