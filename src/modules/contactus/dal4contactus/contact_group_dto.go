package dal4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
)

type ContactGroupDto struct {
	briefs4contactus.ContactGroupBrief
	briefs4contactus.WithMultiTeamContacts[*briefs4contactus.ContactBrief]
}

func (v *ContactGroupDto) Validate() error {
	if err := v.ContactGroupBrief.Validate(); err != nil {
		return err
	}
	if err := v.WithMultiTeamContacts.Validate(); err != nil {
		return err
	}
	return nil
}
