package config

import (
	"net/url"
	"testing"
	"time"

	"golang.org/x/text/encoding"
)

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		Scrapers ScrapersConfig
		Output   OutputConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Scrapers: ScrapersConfig{
					Grob: ScraperConfig{
						BaseURL: url.URL{
							Scheme: "https",
							Host:   "www.gr-oborona.ru",
						},
						Retry: RetryConfig{
							Attempts:       1,
							MinInterval:    1 * time.Second,
							MaxInterval:    1 * time.Second,
							Factor:         1.5,
							MaxElapsedTime: 30 * time.Second,
						},
					},
				},
				Output: OutputConfig{
					FileName: "./out/songs.json",
				},
			},
		},
		{
			name: "err  invalid scrapers",
			fields: fields{
				Scrapers: ScrapersConfig{},
				Output: OutputConfig{
					FileName: "./out/songs.json",
				},
			},
			wantErr: true,
		},
		{
			name: "err  invalid output",
			fields: fields{
				Scrapers: ScrapersConfig{},
				Output: OutputConfig{
					FileName: "./out/songs.json",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Scrapers: tt.fields.Scrapers,
				Output:   tt.fields.Output,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
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
					BaseURL: url.URL{
						Scheme: "https",
						Host:   "www.gr-oborona.ru",
					},
					Retry: RetryConfig{
						Attempts:       1,
						MinInterval:    1 * time.Second,
						MaxInterval:    1 * time.Second,
						Factor:         1.5,
						MaxElapsedTime: 30 * time.Second,
					},
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
		BaseURL    url.URL
		Encoding   encoding.Encoding
		Validation bool
		Retry      RetryConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				BaseURL: url.URL{
					Scheme: "https",
					Host:   "www.gr-oborona.ru",
				},
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
				BaseURL: url.URL{},
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
				BaseURL: url.URL{
					Scheme: "https",
					Host:   "www.gr-oborona.ru",
				},
				Retry: RetryConfig{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ScraperConfig{
				BaseURL:    tt.fields.BaseURL,
				Encoding:   tt.fields.Encoding,
				Validation: tt.fields.Validation,
				Retry:      tt.fields.Retry,
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
