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
	assetEmptyFieldIsNotOmitted("title")
	assetEmptyFieldIsNotOmitted("type")
	assetEmptyFieldIsNotOmitted("status")
	assetEmptyFieldIsNotOmitted("category")
	assetEmptyFieldIsNotOmitted("possession")
	assetEmptyFieldIsNotOmitted("countryID")
	assetEmptyFieldIsNotOmitted("make")
	assetEmptyFieldIsNotOmitted("model")
	assetEmptyFieldIsNotOmitted("regNumber")
}
