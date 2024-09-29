package dal4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dbo4spaceus"
)

type ContactusSpaceEntry = record.DataWithID[string, *dbo4contactus.ContactusSpaceDbo]

func NewContactusSpaceKey(spaceID string) *dal.Key {
	return dbo4spaceus.NewSpaceModuleKey(spaceID, const4contactus.ModuleID)
}

func NewContactusSpaceEntry(spaceID string) ContactusSpaceEntry {
	return NewContactusSpaceEntryWithData(spaceID, new(dbo4contactus.ContactusSpaceDbo))
}

func NewContactusSpaceEntryWithData(spaceID string, data *dbo4contactus.ContactusSpaceDbo) ContactusSpaceEntry {
	return dbo4spaceus.NewSpaceModuleEntry(spaceID, const4contactus.ModuleID, data)
}
