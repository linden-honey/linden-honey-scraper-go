package fetcher

import (
	"net/http"
	"net/url"
	"testing"

	"golang.org/x/text/encoding/charmap"
)

func TestFetcher_Validate(t *testing.T) {
	type fields struct {
		baseURL  *url.URL
		encoding *charmap.Charmap
		client   httpClient
		retry    retryFunc
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
			},
		},
		{
			name: "err  no base url",
			fields: fields{
				baseURL:  nil,
				encoding: charmap.Windows1251,
				client:   &http.Client{},
			},
			wantErr: true,
		},
		{
			name: "err  no encoding",
			fields: fields{
				baseURL:  &url.URL{},
				encoding: nil,
				client:   &http.Client{},
			},
			wantErr: true,
		},
		{
			name: "err  no client",
			fields: fields{
				baseURL:  &url.URL{},
				encoding: charmap.Windows1251,
				client:   nil,
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
