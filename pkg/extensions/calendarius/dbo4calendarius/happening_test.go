package dbo4calendarius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsKnownHappeningStatus(t *testing.T) {
	tests := []struct {
		status string
		want   bool
	}{
		{HappeningStatusActive, true},
		{HappeningStatusArchived, true},
		{HappeningStatusCanceled, true},
		{HappeningStatusDeleted, true},
		{"", false},
		{"unknown", false},
	}
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			assert.Equal(t, tt.want, IsKnownHappeningStatus(tt.status))
		})
	}
}

func TestIsKnownSlotStatus(t *testing.T) {
	assert.True(t, IsKnownSlotStatus(SlotStatusActive))
	assert.True(t, IsKnownSlotStatus(SlotStatusCanceled))
	assert.False(t, IsKnownSlotStatus("unknown"))
}

func TestNewHappeningKey(t *testing.T) {
	key := NewHappeningKey("space1", "h1")
	assert.NotNil(t, key)
	assert.Equal(t, "h1", key.ID)
}

func TestNewHappeningEntry(t *testing.T) {
	entry := NewHappeningEntry("space1", "h1")
	assert.Equal(t, "h1", entry.ID)
	assert.NotNil(t, entry.Key)
	assert.NotNil(t, entry.Data)
	assert.NotNil(t, entry.Record)
}
