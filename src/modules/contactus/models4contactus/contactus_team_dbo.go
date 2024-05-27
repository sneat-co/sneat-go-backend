package models4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
)

type ContactusTeamDbo struct {
	briefs4contactus.WithSingleTeamContactsWithoutContactIDs[*briefs4contactus.ContactBrief]
}

func (v *ContactusTeamDbo) Validate() error {
	return v.WithSingleTeamContactsWithoutContactIDs.Validate()
}
