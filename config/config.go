package config

import (
	"fmt"

	"github.com/linden-honey/linden-honey-sdk-go/env"

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
	Addr string
}

// ScraperConfig represents the scraper configuration
type ScraperConfig struct {
	BaseURL string
}

// NewConfig returns a pointer to the new instance of Config or an error
func NewConfig() (cfg *Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to build a config: %s", r)
		}
	}()

	cfg = &Config{
		Application: ApplicationConfig{
			Scrapers: map[string]ScraperConfig{
				parser.GrobParserID: {
					BaseURL: env.GetEnv(
						"APPLICATION_SCRAPERS_GROB_BASE_URL", "http://www.gr-oborona.ru/",
					),
				},
			},
		},
		Server: ServerConfig{
			Addr: env.GetEnv("SERVER_ADDR", "localhost:8080"),
		},
	}

	return
}
