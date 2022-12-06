package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Config represents a configuration object
type Config struct {
	Server   ServerConfig
	Scrapers ScrapersConfig
}

// ServerConfig represents a configuration object
type ServerConfig struct {
	Host   string       `env:"SERVER_HOST"`
	Port   int          `env:"SERVER_PORT"`
	Health HealthConfig `envPrefix:"SERVER_"`
	Spec   SpecConfig   `envPrefix:"SERVER_"`
}

// HealthConfig represents a configuration object
type HealthConfig struct {
	Enabled bool   `env:"HEALTH_ENABLED"`
	Path    string `env:"HEALTH_PATH"`
}

// SpecConfig represents a configuration object
type SpecConfig struct {
	FilePath string `env:"SPEC_FILE_PATH"`
}

// ScrapersConfig represents a configuration object
type ScrapersConfig struct {
	Grob ScraperConfig `envPrefix:"GROB_"`
}

// ScraperConfig represents a configuration object
type ScraperConfig struct {
	BaseURL string `env:"SCRAPER_BASE_URL"`
}

// New returns a pointer to the new instance of [Config] or an error
func New() (*Config, error) {
	cfg := DefaultConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return &cfg, nil
}
