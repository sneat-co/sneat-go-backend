package dal4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContactusTeamModuleKey(t *testing.T) {
	t.Run("panic_on_empty_team_id", func(t *testing.T) {
		assert.Panics(t, func() {
			NewContactusTeamModuleKey("")
		})
	})
	t.Run("should_pass", func(t *testing.T) {
		entry := NewContactusTeamModuleKey("team1")
		assert.NotNil(t, entry)
		assert.Equal(t, const4contactus.ModuleID, entry.ID)
	})
}

func TestNewContactusTeamModuleEntry(t *testing.T) {
	t.Run("with_id", func(t *testing.T) {
		entry := NewContactusTeamModuleEntry("team1")
		assert.NotNil(t, entry.Data)
		entry.Record.SetError(nil)
		assert.Same(t, entry.Data, entry.Record.Data())
	})
	t.Run("with_empty_id", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = NewContactusTeamModuleEntry("")
		})
	})
}

func TestNewContactusTeamModuleEntryWithData(t *testing.T) {
	t.Run("nil_data", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = NewContactusTeamModuleEntryWithData("team1", nil)
		})
	})
	t.Run("not_nil_data", func(t *testing.T) {
		data := new(models4contactus.ContactusTeamDbo)
		entry := NewContactusTeamModuleEntryWithData("team1", data)
		assert.NotNil(t, entry.Data)
		assert.Same(t, data, entry.Data)
		entry.Record.SetError(nil)
		assert.Same(t, data, entry.Record.Data())
	})
}
