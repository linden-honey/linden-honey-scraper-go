package fetcher

import (
	"net/http"
	"net/url"
	"testing"
)

func TestFetcher_Validate(t *testing.T) {
	type fields struct {
		baseURL url.URL
		client  httpClient
		retry   retryFunc
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				baseURL: url.URL{
					Scheme: "https",
					Host:   "www.google.com",
				},
				client: &http.Client{},
			},
		},
		{
			name: "err  empty url",
			fields: fields{
				baseURL: url.URL{},
				client:  &http.Client{},
			},
			wantErr: true,
		},
		{
			name: "err  no client",
			fields: fields{
				baseURL: url.URL{
					Scheme: "https",
					Host:   "www.google.com",
				},
				client: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher{
				baseURL: tt.fields.baseURL,
				client:  tt.fields.client,
				retry:   tt.fields.retry,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Fetcher.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
