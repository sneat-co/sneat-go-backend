package briefs4assetus

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestAssetBriefNoOmitEmpty(t *testing.T) {
	var assetBrief = AssetBrief{}
	b, err := json.Marshal(assetBrief)
	if err != nil {
		t.Errorf("failed to marshal to JSON: %v", err)
	}
	if b == nil {
		t.Error("marshalled to nil")
	}
	s := string(b)
	assetEmptyFieldIsNotOmitted := func(field string) {
		t.Helper()
		if !strings.Contains(s, fmt.Sprintf(`"%s":`, field)) {
			t.Errorf("should NOT omit empty `%s`", field)
		}
	}
	assetEmptyFieldIsNotOmitted("type")
	assetEmptyFieldIsNotOmitted("status")
	assetEmptyFieldIsNotOmitted("category")
	assetEmptyFieldIsNotOmitted("possession")
	assetEmptyFieldIsNotOmitted("countryID")
}

func TestAssetBrief_Validate(t *testing.T) {
	tests := []struct {
		name    string
		brief   *AssetBrief
		wantErr bool
	}{
		{"nil", nil, true},
		{"empty", &AssetBrief{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.brief.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AssetBrief.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
