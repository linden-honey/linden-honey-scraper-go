package docs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpec_Validate(t *testing.T) {
	tests := []struct {
		name    string
		s       Spec
		wantErr bool
	}{
		{
			name: "valid",
			s:    "some spec",
		},
		{
			name:    "empty phrase",
			s:       "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		rq := require.New(t)

		err := tt.s.Validate()

		if tt.wantErr {
			rq.Error(err)
		} else {
			rq.NoError(err)
		}
	}
}
