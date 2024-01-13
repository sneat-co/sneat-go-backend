package models4brands

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"strings"
	"testing"
)

func TestBrand_Validate(t *testing.T) {
	type fields struct {
		Title      string
		AssetTypes []const4assetus.AssetCategory
		WebsiteURL string
		Models     []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr []string
	}{
		{
			name:    "empty",
			fields:  fields{},
			wantErr: []string{"missing required field", "title"},
		},
		{
			name: "only_title",
			fields: fields{
				Title: "Only title",
			},
			wantErr: []string{"missing required field", "assetTypes"},
		},
		{
			name: "valid",
			fields: fields{
				Title:      "Only title",
				AssetTypes: []string{"cars"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Maker{
				Title:      tt.fields.Title,
				Models:     tt.fields.Models,
				AssetTypes: tt.fields.AssetTypes,
			}
			err := v.Validate()
			if err == nil {
				if len(tt.wantErr) != 0 {
					t.Errorf("Validation passed while expected to get error(s): %+v", tt.wantErr)
				}
			} else /* if err != nil */ {
				if len(tt.wantErr) == 0 {
					t.Errorf("Validate() returned unepxcted error: %v", err)
				} else {
					for _, wantErr := range tt.wantErr {
						if !strings.Contains(err.Error(), wantErr) {
							t.Errorf("Validate() expected to return an error with %v but returned : %v", tt.wantErr, err)
							continue
						}
					}
				}
			}
		})
	}
}
