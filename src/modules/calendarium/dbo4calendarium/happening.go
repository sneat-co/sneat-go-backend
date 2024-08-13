package dbo4calendarium

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
)

//const SingleHappeningsCollection = "single_happenings"

// NewHappeningKey creates a new happening key
func NewHappeningKey(teamID, happeningID string) *dal.Key {
	return dbo4spaceus.NewSpaceModuleItemKey(teamID, const4calendarium.ModuleID, const4calendarium.HappeningsCollection, happeningID)
}

// HappeningType is either "recurring" or "single"
type HappeningType = string

const (
	// HappeningTypeRecurring = "recurring"
	HappeningTypeRecurring HappeningType = "recurring"

	// HappeningTypeSingle = "single"
	HappeningTypeSingle HappeningType = "single"
)

const (
	HappeningStatusActive   = "active"
	HappeningStatusArchived = "archived"
	HappeningStatusCanceled = "canceled"
	HappeningStatusDeleted  = "deleted"
)

// IsKnownHappeningStatus detects if a string is a know happening status
func IsKnownHappeningStatus(status string) bool {
	switch status {
	case HappeningStatusActive, HappeningStatusArchived, HappeningStatusCanceled, HappeningStatusDeleted:
		return true
	}
	return false
}

func NewHappeningEntry(teamID, happeningID string) (v HappeningEntry) {
	return NewHappeningEntryWithDbo(teamID, happeningID, new(HappeningDbo))
}

func NewHappeningEntryWithDbo(teamID, happeningID string, dto *HappeningDbo) (v HappeningEntry) {
	v.ID = happeningID
	v.Key = NewHappeningKey(teamID, happeningID)
	v.Data = dto
	v.Record = dal.NewRecordWithData(v.Key, dto)
	return
}
