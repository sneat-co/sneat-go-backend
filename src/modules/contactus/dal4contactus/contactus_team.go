package dal4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
)

type ContactusTeamModuleEntry = record.DataWithID[string, *models4contactus.ContactusTeamDbo]

func NewContactusTeamModuleKey(teamID string) *dal.Key {
	return dal4teamus.NewTeamModuleKey(teamID, const4contactus.ModuleID)
}

func NewContactusTeamModuleEntry(teamID string) ContactusTeamModuleEntry {
	return NewContactusTeamModuleEntryWithData(teamID, new(models4contactus.ContactusTeamDbo))
}

func NewContactusTeamModuleEntryWithData(teamID string, data *models4contactus.ContactusTeamDbo) ContactusTeamModuleEntry {
	return dal4teamus.NewTeamModuleEntry(teamID, const4contactus.ModuleID, data)
}
