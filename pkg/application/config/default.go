package config

import (
	"time"
)

func Default() Config {
	return Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Health: HealthConfig{
			Enabled: true,
			Path:    "/health",
		},
		Spec: SpecConfig{
			FilePath: "./api/openapi.json",
		},
		Scrapers: ScrapersConfig{
			Grob: ScraperConfig{
				BaseURL: "https://www.gr-oborona.ru/",
				Retry: RetryConfig{
					Attempts:       5,
					MinInterval:    2 * time.Second,
					MaxInterval:    10 * time.Second,
					Factor:         1.5,
					MaxElapsedTime: 30 * time.Second,
				},
			},
		},
	}
}
