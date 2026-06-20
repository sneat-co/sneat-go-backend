package dbo4companius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompanyBase_Validate_more(t *testing.T) {
	tests := []struct {
		name    string
		base    CompanyBase
		wantErr bool
	}{
		{"valid_work", CompanyBase{Kind: "work", Type: "work", Title: "Acme"}, false},
		{"missing_kind", CompanyBase{Type: "work", Title: "Acme"}, true},
		{"missing_type", CompanyBase{Kind: "work", Title: "Acme"}, true},
		{"work_missing_title", CompanyBase{Kind: "work", Type: "work"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.base.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCompanyDbo_Validate(t *testing.T) {
	validBase := CompanyBase{Kind: "work", Type: "work", Title: "Acme"}
	tests := []struct {
		name    string
		dbo     CompanyDbo
		wantErr bool
	}{
		{"valid_no_numbers", CompanyDbo{CompanyBase: validBase}, false},
		{"valid_with_numbers", CompanyDbo{
			CompanyBase: validBase,
			NumberOf:    map[string]int{"members": 3, "documents": 0},
		}, false},
		{"invalid_base", CompanyDbo{CompanyBase: CompanyBase{}}, true},
		{"negative_number", CompanyDbo{
			CompanyBase: validBase,
			NumberOf:    map[string]int{"members": -1},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dbo.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
