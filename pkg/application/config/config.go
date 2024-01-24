package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

// TODO: mark required fields like this `env:"ENV_NAME,required"`
// TODO: use `envPrefix:"PREFIX_"` and default env names

// Config is a configuration object.
type Config struct {
	Server   ServerConfig
	Health   HealthConfig
	Spec     SpecConfig
	Scrapers ScrapersConfig
}

// ServerConfig is a configuration object.
type ServerConfig struct {
	Host string `env:"SERVER_HOST"`
	Port int    `env:"SERVER_PORT"`
}

// HealthConfig is a configuration object.
type HealthConfig struct {
	Enabled bool   `env:"HEALTH_ENABLED"`
	Path    string `env:"HEALTH_PATH"`
}

// SpecConfig is a configuration object.
type SpecConfig struct {
	FilePath string `env:"SPEC_FILE_PATH"`
}

// ScrapersConfig is a configuration object.
type ScrapersConfig struct {
	Grob ScraperConfig `envPrefix:"GROB_"`
}

// ScraperConfig is a configuration object.
type ScraperConfig struct {
	BaseURL string `env:"SCRAPER_BASE_URL"`
	Retry   RetryConfig
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

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate a config: %w", err)
	}

	return &cfg, nil
}
