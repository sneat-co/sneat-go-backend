package dbo4sportus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpot_Validate(t *testing.T) {
	tests := []struct {
		name    string
		spot    Spot
		wantErr bool
	}{
		{"valid", Spot{SpotBrief: SpotBrief{Title: "Spot 1"}}, false},
		{"empty_title", Spot{SpotBrief: SpotBrief{Title: ""}}, true},
		{"whitespace_title", Spot{SpotBrief: SpotBrief{Title: "  "}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spot.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
