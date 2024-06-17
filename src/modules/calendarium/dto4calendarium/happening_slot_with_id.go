package dto4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/strongo/validation"
)

type HappeningSlotWithID struct {
	ID string `json:"id"`
	dbo4calendarium.HappeningSlot
}

func (v HappeningSlotWithID) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	return v.HappeningSlot.Validate()
}
