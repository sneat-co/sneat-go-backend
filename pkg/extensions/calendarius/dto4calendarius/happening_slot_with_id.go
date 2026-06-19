package dto4calendarius

import (
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/strongo/validation"
)

type HappeningSlotWithID struct {
	ID string `json:"id"`
	dbo4calendarius.HappeningSlot
}

func (v HappeningSlotWithID) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	return v.HappeningSlot.Validate()
}
