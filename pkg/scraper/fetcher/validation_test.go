package fetcher

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func TestFetcher_Validate(t *testing.T) {
	type fields struct {
		baseURL  *url.URL
		encoding *charmap.Charmap
		client   httpClient
		retry    *RetryConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				baseURL:  &url.URL{},
				encoding: charmap.Windows1251,
				client:   &http.Client{},
				retry: &RetryConfig{
					Attempts:   3,
					MinTimeout: 1 * time.Second,
					MaxTimeout: 6 * time.Second,
					Factor:     3 * time.Second,
				},
			},
		},
		{
			name: "err  no base url",
			fields: fields{
				baseURL:  nil,
				encoding: charmap.Windows1251,
				client:   &http.Client{},
				retry: &RetryConfig{
					Attempts:   3,
					MinTimeout: 1 * time.Second,
					MaxTimeout: 6 * time.Second,
					Factor:     3 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "err  no encoding",
			fields: fields{
				baseURL:  &url.URL{},
				encoding: nil,
				client:   &http.Client{},
				retry: &RetryConfig{
					Attempts:   3,
					MinTimeout: 1 * time.Second,
					MaxTimeout: 6 * time.Second,
					Factor:     3 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "err  no client",
			fields: fields{
				baseURL:  &url.URL{},
				encoding: charmap.Windows1251,
				client:   nil,
				retry: &RetryConfig{
					Attempts:   3,
					MinTimeout: 1 * time.Second,
					MaxTimeout: 6 * time.Second,
					Factor:     3 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "err  invalid retry",
			fields: fields{
				baseURL:  &url.URL{},
				encoding: charmap.Windows1251,
				client:   &http.Client{},
				retry:    &RetryConfig{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher{
				baseURL:  tt.fields.baseURL,
				encoding: tt.fields.encoding,
				client:   tt.fields.client,
				retry:    tt.fields.retry,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Fetcher.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRetryConfig_Validate(t *testing.T) {
	type fields struct {
		Attempts   int
		MinTimeout time.Duration
		MaxTimeout time.Duration
		Factor     time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Attempts:   3,
				MinTimeout: 1 * time.Second,
				MaxTimeout: 6 * time.Second,
				Factor:     3 * time.Second,
			},
		},
		{
			name: "err  attempts is non-positive number",
			fields: fields{
				Attempts:   0,
				MinTimeout: 1 * time.Second,
				MaxTimeout: 6 * time.Second,
				Factor:     3 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "err  factor is non-positive number",
			fields: fields{
				Attempts:   3,
				MinTimeout: 1 * time.Second,
				MaxTimeout: 6 * time.Second,
				Factor:     0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := RetryConfig{
				Attempts:   tt.fields.Attempts,
				MinTimeout: tt.fields.MinTimeout,
				MaxTimeout: tt.fields.MaxTimeout,
				Factor:     tt.fields.Factor,
			}
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("RetryConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
