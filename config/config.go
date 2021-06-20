package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/parser"
)

// Config represents the root configuration
type Config struct {
	Application ApplicationConfig
	Server      ServerConfig
}

// ApplicationConfig represents the application configuration
type ApplicationConfig struct {
	Scrapers map[string]ScraperConfig
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Addr string `env:"SERVER_ADDR"`
}

// ScraperConfig represents the scraper configuration
type ScraperConfig struct {
	BaseURL string
}

// NewConfig returns a pointer to the new instance of Config or an error
func NewConfig() (*Config, error) {
	cfg := &Config{
		Application: ApplicationConfig{
			Scrapers: map[string]ScraperConfig{
				parser.GrobParserID: {
					BaseURL: "http://www.gr-oborona.ru/",
				},
			},
		},
		Server: ServerConfig{
			Addr: "localhost:8080",
		},
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
