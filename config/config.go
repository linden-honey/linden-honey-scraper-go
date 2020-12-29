package config

import (
	"fmt"

	"github.com/spf13/cast"

	"github.com/linden-honey/linden-honey-sdk-go/env"
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
	Host string
	Port int
}

// ScraperConfig represents the scraper configuration
type ScraperConfig struct {
	BaseURL string
}

// NewConfig returns a pointer to the new instance of Config or an error
func NewConfig() (cfg *Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to build config: %s", r)
		}
	}()

	cfg = &Config{
		Application: ApplicationConfig{
			Scrapers: map[string]ScraperConfig{
				"grob": {
					BaseURL: env.GetEnv(
						"APPLICATION_SCRAPERS_GROB_BASE_URL", "http://www.gr-oborona.ru/",
					),
				},
			},
		},
		Server: ServerConfig{
			Host: env.GetEnv("SERVER_HOST", "0.0.0.0"),
			Port: cast.ToInt(env.GetEnv("SERVER_PORT", "8080")),
		},
	}

	return
}
