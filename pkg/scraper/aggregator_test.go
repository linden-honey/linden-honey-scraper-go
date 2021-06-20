package scraper

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_aggregationErr_Error(t *testing.T) {
	type fields struct {
		msg     string
		reasons []error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				msg: "failed to aggregate scraped songs",
				reasons: []error{
					fmt.Errorf(
						"failed to get song with id 1056901568: %w",
						errors.New("song is invalid"),
					),
				},
			},
			want: "failed to aggregate scraped songs: [failed to get song with id 1056901568: song is invalid]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rq := require.New(t)

			err := &aggregationErr{
				msg:     tt.fields.msg,
				reasons: tt.fields.reasons,
			}
			got := err.Error()

			rq.Equal(tt.want, got)
		})
	}
}
