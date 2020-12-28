package config

import (
	"fmt"

	"github.com/spf13/cast"

	"github.com/linden-honey/linden-honey-sdk-go/env"
)

type Config struct {
	Application ApplicationConfig
	Server      ServerConfig
}

type ApplicationConfig struct {
	Scrapers map[string]ScraperConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type ScraperConfig struct {
	BaseURL string
}

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
			Host: env.GetEnv("SERVER_HOST", "127.0.0.1"),
			Port: cast.ToInt(env.GetEnv("SERVER_PORT", "8080")),
		},
	}

	return
}
