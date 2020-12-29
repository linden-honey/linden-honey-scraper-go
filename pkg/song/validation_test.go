package song

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuote_Validate(t *testing.T) {
	type fields struct {
		Phrase string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Phrase: "Some phrase",
			},
		},
		{
			name: "empty phrase",
			fields: fields{
				Phrase: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := require.New(t)

			q := Quote{
				Phrase: tt.fields.Phrase,
			}
			err := q.Validate()

			if tt.wantErr {
				rq.Error(err)
			} else {
				rq.NoError(err)
			}
		})
	}
}

func TestVerse_Validate(t *testing.T) {
	type fields struct {
		Quotes []Quote
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Quotes: []Quote{
					{
						Phrase: "Some phrase",
					},
				},
			},
		},
		{
			name: "empty quotes",
			fields: fields{
				Quotes: make([]Quote, 0),
			},
			wantErr: true,
		},
		{
			name: "invalid quote",
			fields: fields{
				Quotes: []Quote{
					{
						Phrase: "Some phrase",
					},
					{
						Phrase: "",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := require.New(t)

			v := Verse{
				Quotes: tt.fields.Quotes,
			}
			err := v.Validate()

			if tt.wantErr {
				rq.Error(err)
			} else {
				rq.NoError(err)
			}
		})
	}
}

func TestSong_Validate(t *testing.T) {
	type fields struct {
		Title  string
		Author string
		Album  string
		Verses []Verse
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Title:  "Some title",
				Author: "Some author",
				Album:  "Some album",
				Verses: []Verse{
					{
						Quotes: []Quote{
							{
								Phrase: "Some phrase",
							},
						},
					},
				},
			},
		},
		{
			name: "empty title",
			fields: fields{
				Title:  "",
				Author: "Some author",
				Album:  "Some album",
				Verses: []Verse{
					{
						Quotes: []Quote{
							{
								Phrase: "Some phrase",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty verses",
			fields: fields{
				Title:  "Some title",
				Author: "Some author",
				Album:  "Some album",
				Verses: make([]Verse, 0),
			},
			wantErr: true,
		},
		{
			name: "invalid verse",
			fields: fields{
				Title:  "Some title",
				Author: "Some author",
				Album:  "Some album",
				Verses: []Verse{
					{
						Quotes: []Quote{
							{
								Phrase: "Some phrase",
							},
						},
					},
					{
						Quotes: make([]Quote, 0),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid quote",
			fields: fields{
				Title:  "Some title",
				Author: "Some author",
				Album:  "Some album",
				Verses: []Verse{
					{
						Quotes: []Quote{
							{
								Phrase: "Some phrase",
							},
						},
					},
					{
						Quotes: []Quote{
							{
								Phrase: "",
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := require.New(t)

			s := Song{
				Title:  tt.fields.Title,
				Author: tt.fields.Author,
				Album:  tt.fields.Album,
				Verses: tt.fields.Verses,
			}
			err := s.Validate()

			if tt.wantErr {
				rq.Error(err)
			} else {
				rq.NoError(err)
			}
		})
	}
}

func TestPreview_Validate(t *testing.T) {
	type fields struct {
		ID    string
		Title string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				ID:    "123",
				Title: "Some title",
			},
		},
		{
			name: "empty id",
			fields: fields{
				ID:    "",
				Title: "Some title",
			},
			wantErr: true,
		},
		{
			name: "empty title",
			fields: fields{
				ID:    "123",
				Title: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := require.New(t)

			p := Preview{
				ID:    tt.fields.ID,
				Title: tt.fields.Title,
			}
			err := p.Validate()

			if tt.wantErr {
				rq.Error(err)
			} else {
				rq.NoError(err)
			}
		})
	}
}
