package config

import (
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		Server   ServerConfig
		Health   HealthConfig
		Spec     SpecConfig
		Scrapers ScrapersConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
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
						BaseURL: "https://test.com/",
					},
				},
			},
		},
		{
			name: "err  invalid server",
			fields: fields{
				Server: ServerConfig{},
				Scrapers: ScrapersConfig{
					Grob: ScraperConfig{
						BaseURL: "https://test.com/",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "err  invalid health",
			fields: fields{
				Server: ServerConfig{
					Host: "localhost",
					Port: 8080,
				},
				Health: HealthConfig{},
				Spec: SpecConfig{
					FilePath: "./api/openapi.json",
				},
				Scrapers: ScrapersConfig{
					Grob: ScraperConfig{
						BaseURL: "https://test.com/",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "err  invalid spec",
			fields: fields{
				Server: ServerConfig{
					Host: "localhost",
					Port: 8080,
				},
				Health: HealthConfig{
					Enabled: true,
					Path:    "/health",
				},
				Spec: SpecConfig{},
				Scrapers: ScrapersConfig{
					Grob: ScraperConfig{
						BaseURL: "https://test.com/",
					},
				},
			},
			wantErr: true,
		},

		{
			name: "err  invalid scrapers",
			fields: fields{
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
				Scrapers: ScrapersConfig{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Server:   tt.fields.Server,
				Scrapers: tt.fields.Scrapers,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerConfig_Validate(t *testing.T) {
	type fields struct {
		Host string
		Port int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Host: "localhost",
				Port: 8080,
			},
		},
		{
			name: "err  empty host",
			fields: fields{
				Host: "",
				Port: 8080,
			},
			wantErr: true,
		},
		{
			name: "err  invalid port",
			fields: fields{
				Host: "localhost",
				Port: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ServerConfig{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ServerConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHealthConfig_Validate(t *testing.T) {
	type fields struct {
		Enabled bool
		Path    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Enabled: true,
				Path:    "/health",
			},
		},
		{
			name: "err  empty path",
			fields: fields{
				Enabled: true,
				Path:    "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := HealthConfig{
				Enabled: tt.fields.Enabled,
				Path:    tt.fields.Path,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("HealthConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSpecConfig_Validate(t *testing.T) {
	type fields struct {
		FilePath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				FilePath: "./api/openapi.json",
			},
		},
		{
			name: "err  empty file path",
			fields: fields{
				FilePath: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := SpecConfig{
				FilePath: tt.fields.FilePath,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SpecConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScrapersConfig_Validate(t *testing.T) {
	type fields struct {
		Grob ScraperConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Grob: ScraperConfig{
					BaseURL: "https://test.com/",
				},
			},
		},
		{
			name: "err  invalid grob scraper config",
			fields: fields{
				Grob: ScraperConfig{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ScrapersConfig{
				Grob: tt.fields.Grob,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ScrapersConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScraperConfig_Validate(t *testing.T) {
	type fields struct {
		BaseURL string
		Retry   RetryConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				BaseURL: "https://test.com/",
				Retry: RetryConfig{
					Attempts:       1,
					MinInterval:    1 * time.Second,
					MaxInterval:    1 * time.Second,
					Factor:         1.5,
					MaxElapsedTime: 30 * time.Second,
				},
			},
		},
		{
			name: "err  empty base url",
			fields: fields{
				BaseURL: "",
				Retry: RetryConfig{
					Attempts:       1,
					MinInterval:    1 * time.Second,
					MaxInterval:    1 * time.Second,
					Factor:         1.5,
					MaxElapsedTime: 30 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "err  invalid retry config",
			fields: fields{
				BaseURL: "https://test.com/",
				Retry:   RetryConfig{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ScraperConfig{
				BaseURL: tt.fields.BaseURL,
				Retry:   tt.fields.Retry,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ScraperConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRetryConfig_Validate(t *testing.T) {
	type fields struct {
		Attempts       uint
		MinInterval    time.Duration
		MaxInterval    time.Duration
		Factor         float64
		MaxElapsedTime time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Attempts:       3,
				MinInterval:    1 * time.Second,
				MaxInterval:    10 * time.Second,
				Factor:         1.5,
				MaxElapsedTime: 30 * time.Second,
			},
		},
		{
			name: "err  attempts is non-positive number",
			fields: fields{
				Attempts:       0,
				MinInterval:    1 * time.Second,
				MaxInterval:    6 * time.Second,
				Factor:         1.5,
				MaxElapsedTime: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "err  min interval is non-positive number",
			fields: fields{
				Attempts:       3,
				MinInterval:    0,
				MaxInterval:    6 * time.Second,
				Factor:         1.5,
				MaxElapsedTime: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "err  max interval is non-positive number",
			fields: fields{
				Attempts:       3,
				MinInterval:    1 * time.Second,
				MaxInterval:    0,
				Factor:         1.5,
				MaxElapsedTime: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "err  min interval is greater than max interval",
			fields: fields{
				Attempts:    0,
				MinInterval: 1 * time.Second,
				MaxInterval: 6 * time.Second,
				Factor:      1.5,
			},
			wantErr: true,
		},
		{
			name: "err  factor is non-positive number",
			fields: fields{
				Attempts:       3,
				MinInterval:    1 * time.Second,
				MaxInterval:    6 * time.Second,
				Factor:         0,
				MaxElapsedTime: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "err  max elapsed time is lower than min interval",
			fields: fields{
				Attempts:       3,
				MinInterval:    1 * time.Second,
				MaxInterval:    6 * time.Second,
				Factor:         0,
				MaxElapsedTime: 30 * time.Second,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := RetryConfig{
				Attempts:    tt.fields.Attempts,
				MinInterval: tt.fields.MinInterval,
				MaxInterval: tt.fields.MaxInterval,
				Factor:      tt.fields.Factor,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("RetryConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
