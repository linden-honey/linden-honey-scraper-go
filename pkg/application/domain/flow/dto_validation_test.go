package flow

import "testing"

func TestSimpleFlowInput_Validate(t *testing.T) {
	type fields struct {
		OutputFileName string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				OutputFileName: "./out/songs.json",
			},
		},
		{
			name: "err  empty file name",
			fields: fields{
				OutputFileName: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := SimpleFlowInput{
				OutputFileName: tt.fields.OutputFileName,
			}
			if err := dto.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SimpleFlowInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
