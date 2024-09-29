package dbo4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/briefs4contactus"
)

type ContactusSpaceDbo struct {
	TotalContactsCountByStatus map[string]int `json:"totalContactsCountByStatus,omitempty" firestore:"totalContactsCountByStatus,omitempty"`
	briefs4contactus.WithSingleSpaceContactsWithoutContactIDs[*briefs4contactus.ContactBrief]
}

func (v *ContactusSpaceDbo) Validate() error {
	return v.WithSingleSpaceContactsWithoutContactIDs.Validate()
}
