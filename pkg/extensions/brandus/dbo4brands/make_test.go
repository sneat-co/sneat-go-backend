package dbo4brands

import (
	"strings"
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/const4assetus"
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
		{
			name: "valid_with_url",
			fields: fields{
				Title:      "Toyota",
				AssetTypes: []string{"cars"},
				WebsiteURL: "https://toyota.com",
			},
		},
		{
			name: "invalid_url",
			fields: fields{
				Title:      "Toyota",
				AssetTypes: []string{"cars"},
				WebsiteURL: " : invalid",
			},
			wantErr: []string{"bad value for field [websiteURL]"},
		},
		{
			name: "invalid_model",
			fields: fields{
				Title:      "Toyota",
				AssetTypes: []string{"cars"},
				Models:     []string{""},
			},
			wantErr: []string{"bad value for field [models[0]]", "missing required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Maker{
				Title:      tt.fields.Title,
				Models:     tt.fields.Models,
				AssetTypes: tt.fields.AssetTypes,
				WebsiteURL: tt.fields.WebsiteURL,
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
