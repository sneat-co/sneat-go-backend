package dto4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
)

type HappeningContactRequest struct {
	HappeningRequest
	Contact models4linkage.ShortTeamModuleDocRef `json:"contact"`
}

func (v HappeningContactRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if err := v.Contact.Validate(); err != nil {
		return err
	}
	return nil
}
