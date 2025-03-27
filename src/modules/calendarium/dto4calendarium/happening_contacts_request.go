package dto4calendarium

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/strongo/validation"
)

type HappeningContactsRequest struct {
	HappeningRequest
	Contacts []dbo4linkage.ShortSpaceModuleDocRef `json:"contacts"`
}

func (v HappeningContactsRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	for i, contact := range v.Contacts {
		if err := contact.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contacts[%d]", i), err.Error())
		}
	}
	return nil
}
