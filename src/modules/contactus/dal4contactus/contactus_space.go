package dal4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
)

type ContactusSpaceModuleEntry = record.DataWithID[string, *models4contactus.ContactusSpaceDbo]

func NewContactusSpaceModuleKey(spaceID string) *dal.Key {
	return dal4spaceus.NewSpaceModuleKey(spaceID, const4contactus.ModuleID)
}

func NewContactusSpaceModuleEntry(spaceID string) ContactusSpaceModuleEntry {
	return NewContactusSpaceModuleEntryWithData(spaceID, new(models4contactus.ContactusSpaceDbo))
}

func NewContactusSpaceModuleEntryWithData(spaceID string, data *models4contactus.ContactusSpaceDbo) ContactusSpaceModuleEntry {
	return dal4spaceus.NewSpaceModuleEntry(spaceID, const4contactus.ModuleID, data)
}
