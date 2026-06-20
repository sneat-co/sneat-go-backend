package dbo4sportus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestItem_Validate(t *testing.T) {
	created := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := created.Add(time.Hour)

	tests := []struct {
		name    string
		item    Item
		wantErr bool
	}{
		{"valid", Item{UserID: "u1", Status: "active", DtCreated: created, DtUpdated: updated}, false},
		{"valid_deleted", Item{UserID: "u1", Status: "deleted", DtCreated: created, DtUpdated: updated}, false},
		{"missing_user", Item{Status: "active", DtCreated: created, DtUpdated: updated}, true},
		{"whitespace_user", Item{UserID: "  ", Status: "active", DtCreated: created, DtUpdated: updated}, true},
		{"missing_dtCreated", Item{UserID: "u1", Status: "active", DtUpdated: updated}, true},
		{"missing_dtUpdated", Item{UserID: "u1", Status: "active", DtCreated: created}, true},
		{"updated_before_created", Item{UserID: "u1", Status: "active", DtCreated: updated, DtUpdated: created}, true},
		{"status_with_spaces", Item{UserID: "u1", Status: " active", DtCreated: created, DtUpdated: updated}, true},
		{"missing_status", Item{UserID: "u1", Status: "", DtCreated: created, DtUpdated: updated}, true},
		{"unknown_status", Item{UserID: "u1", Status: "archived", DtCreated: created, DtUpdated: updated}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
