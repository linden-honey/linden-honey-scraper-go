package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Config represents a configuration object
type Config struct {
	Server   ServerConfig
	Health   HealthConfig
	Scrapers ScraperConfigs
}

// ServerConfig represents a configuration object
type ServerConfig struct {
	Host string `env:"SERVER_HOST"`
	Port int    `env:"SERVER_PORT"`
}

// HealthConfig represents a configuration object
type HealthConfig struct {
	Enabled bool   `env:"HEALTH_ENABLED"`
	Path    string `env:"HEALTH_PATH"`
}

// ScraperConfig represents a configuration object
type ScraperConfig struct {
	BaseURL string `env:"SCRAPER_BASE_URL"`
}

type ScraperConfigs struct {
	Grob ScraperConfig `envPrefix:"GROB_"`
}

// NewConfig returns a pointer to the new instance of Config or an error
func NewConfig() (*Config, error) {
	cfg := DefaultConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return &cfg, nil
}
