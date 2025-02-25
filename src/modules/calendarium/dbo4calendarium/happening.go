package dbo4calendarium

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

//const SingleHappeningsCollection = "single_happenings"

// NewHappeningKey creates a new happening key
func NewHappeningKey(spaceID coretypes.SpaceID, happeningID string) *dal.Key {
	return dbo4spaceus.NewSpaceModuleItemKey(spaceID, const4calendarium.ModuleID, const4calendarium.HappeningsCollection, happeningID)
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

func NewHappeningEntry(spaceID coretypes.SpaceID, happeningID string) (v HappeningEntry) {
	return NewHappeningEntryWithDbo(spaceID, happeningID, new(HappeningDbo))
}

func NewHappeningEntryWithDbo(spaceID coretypes.SpaceID, happeningID string, dto *HappeningDbo) (v HappeningEntry) {
	v.ID = happeningID
	v.Key = NewHappeningKey(spaceID, happeningID)
	v.Data = dto
	v.Record = dal.NewRecordWithData(v.Key, dto)
	return
}
