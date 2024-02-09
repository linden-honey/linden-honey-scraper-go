package flow

import "testing"

func TestRunSimpleFlowRequest_Validate(t *testing.T) {
	type fields struct {
		ArtifactName string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				ArtifactName: "./out/songs.json",
			},
		},
		{
			name: "err  empty file name",
			fields: fields{
				ArtifactName: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := RunSimpleFlowRequest{
				ArtifactName: tt.fields.ArtifactName,
			}
			if err := dto.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SimpleFlowInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
