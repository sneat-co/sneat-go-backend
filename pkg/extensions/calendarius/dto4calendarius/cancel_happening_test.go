package dto4calendarius

import (
	"strings"
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/stretchr/testify/assert"
)

func TestCancelHappeningRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CancelHappeningRequest
		wantErr bool
	}{
		{"valid_no_date", CancelHappeningRequest{
			HappeningRequest: validHappeningRequest(),
		}, false},
		{"valid_with_date_and_slot", CancelHappeningRequest{
			HappeningRequest: validHappeningRequest(),
			Date:             "2020-01-01",
			SlotID:           "slot1",
		}, false},
		{"valid_with_reason", CancelHappeningRequest{
			HappeningRequest: validHappeningRequest(),
			Reason:           "some reason",
		}, false},
		{"invalid_happening_request", CancelHappeningRequest{
			HappeningRequest: HappeningRequest{},
		}, true},
		{"bad_date", CancelHappeningRequest{
			HappeningRequest: validHappeningRequest(),
			Date:             "not-a-date",
			SlotID:           "slot1",
		}, true},
		{"date_without_slot", CancelHappeningRequest{
			HappeningRequest: validHappeningRequest(),
			Date:             "2020-01-01",
		}, true},
		{"reason_too_long", CancelHappeningRequest{
			HappeningRequest: validHappeningRequest(),
			Reason:           strings.Repeat("a", dbo4calendarius.ReasonMaxLen+1),
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
