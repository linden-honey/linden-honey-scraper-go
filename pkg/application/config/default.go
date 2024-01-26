package config

import (
	"net/url"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func Default() Config {
	return Config{
		Scrapers: ScrapersConfig{
			Grob: ScraperConfig{
				BaseURL: url.URL{
					Scheme: "https",
					Host:   "www.gr-oborona.ru",
				},
				Encoding: charmap.Windows1251,
				Retry: RetryConfig{
					Enabled:        true,
					Attempts:       5,
					MinInterval:    2 * time.Second,
					MaxInterval:    10 * time.Second,
					Factor:         1.5,
					MaxElapsedTime: 30 * time.Second,
				},
			},
		},
		Output: OutputConfig{
			FileName: "./songs.json",
		},
	}
}
