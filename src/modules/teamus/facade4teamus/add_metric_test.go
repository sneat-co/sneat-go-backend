package facade4teamus

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAddTeamMetricRequest_Decoding(t *testing.T) {
	decoder := json.NewDecoder(strings.NewReader(`{"metric":{"title": "Test metric"}}`))
	var request AddTeamMetricRequest
	if err := decoder.Decode(&request); err != nil {
		t.Fatal(err)
	}
	if request.Metric.Title != "Test metric" {
		t.Errorf("Title expected to be [Test metric], got: %v", request.Metric.Title)
	}
}
