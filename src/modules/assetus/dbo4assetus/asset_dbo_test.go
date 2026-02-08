package dbo4assetus

import "testing"

func TestAssetDates_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dates   AssetDates
		wantErr bool
	}{
		{"empty", AssetDates{}, false},
		{"valid_dates", AssetDates{
			DateOfBuild:       "2020-01-01",
			DateOfPurchase:    "2020-02-01",
			DateInsuredTill:   "2021-01-01",
			DateCertifiedTill: "2021-02-01",
		}, false},
		{"invalid_build_date", AssetDates{DateOfBuild: "invalid"}, true},
		{"invalid_purchase_date", AssetDates{DateOfPurchase: "invalid"}, true},
		{"invalid_insured_date", AssetDates{DateInsuredTill: "invalid"}, true},
		{"invalid_certified_date", AssetDates{DateCertifiedTill: "invalid"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dates.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AssetDates.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
