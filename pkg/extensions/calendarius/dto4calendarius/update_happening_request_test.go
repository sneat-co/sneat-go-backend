package dto4calendarius

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateHappeningRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     UpdateHappeningRequest
		wantErr bool
	}{
		{"valid", UpdateHappeningRequest{
			HappeningRequest: validHappeningRequest(),
			Title:            "New title",
		}, false},
		{"valid_with_summary_and_description", UpdateHappeningRequest{
			HappeningRequest: validHappeningRequest(),
			Title:            "New title",
			Summary:          "A summary",
			Description:      "A description",
		}, false},
		{"invalid_happening_request", UpdateHappeningRequest{
			HappeningRequest: HappeningRequest{},
			Title:            "New title",
		}, true},
		{"missing_title", UpdateHappeningRequest{
			HappeningRequest: validHappeningRequest(),
		}, true},
		{"title_too_long", UpdateHappeningRequest{
			HappeningRequest: validHappeningRequest(),
			Title:            strings.Repeat("a", 101),
		}, true},
		{"summary_too_long", UpdateHappeningRequest{
			HappeningRequest: validHappeningRequest(),
			Title:            "New title",
			Summary:          strings.Repeat("a", 201),
		}, true},
		{"description_too_long", UpdateHappeningRequest{
			HappeningRequest: validHappeningRequest(),
			Title:            "New title",
			Description:      strings.Repeat("a", 5001),
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
