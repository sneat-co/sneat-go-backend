package facade4spaceus

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAddSpaceMetricRequest_Decoding(t *testing.T) {
	decoder := json.NewDecoder(strings.NewReader(`{"metric":{"title": "Test metric"}}`))
	var request AddSpaceMetricRequest
	if err := decoder.Decode(&request); err != nil {
		t.Fatal(err)
	}
	if request.Metric.Title != "Test metric" {
		t.Errorf("Title expected to be [Test metric], got: %s", request.Metric.Title)
	}
}
