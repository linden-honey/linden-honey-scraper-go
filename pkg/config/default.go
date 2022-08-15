package config

var (
	DefaultConfig = Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Health: HealthConfig{
			Enabled: true,
			Path:    "/health",
		},
		Scrapers: ScraperConfigs{
			"grob": ScraperConfig{
				BaseURL: "http://www.gr-oborona.ru/",
			},
		},
	}
)
