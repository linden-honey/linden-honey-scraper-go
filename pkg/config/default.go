package config

var (
	DefaultConfig = Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
			Health: HealthConfig{
				Enabled: true,
				Path:    "/health",
			},
		},
		Scrapers: ScrapersConfig{
			Grob: ScraperConfig{
				BaseURL: "https://www.gr-oborona.ru/",
			},
		},
	}
)
