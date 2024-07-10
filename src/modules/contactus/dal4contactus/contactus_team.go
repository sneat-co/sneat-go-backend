package dal4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
)

type ContactusSpaceModuleEntry = record.DataWithID[string, *models4contactus.ContactusSpaceDbo]

func NewContactusSpaceModuleKey(teamID string) *dal.Key {
	return dal4teamus.NewSpaceModuleKey(teamID, const4contactus.ModuleID)
}

func NewContactusSpaceModuleEntry(teamID string) ContactusSpaceModuleEntry {
	return NewContactusSpaceModuleEntryWithData(teamID, new(models4contactus.ContactusSpaceDbo))
}

func NewContactusSpaceModuleEntryWithData(teamID string, data *models4contactus.ContactusSpaceDbo) ContactusSpaceModuleEntry {
	return dal4teamus.NewSpaceModuleEntry(teamID, const4contactus.ModuleID, data)
}
