package dbo4retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dbo4teamus"
	"github.com/strongo/validation"
)

type RetroSpaceDbo struct {
	RetroSettings `json:"retroSettings,omitempty" firestore:"retroSettings,omitempty"`
	UpcomingRetro *RetrospectiveCounts         `json:"upcomingRetro,omitempty" firestore:"upcomingRetro,omitempty"`
	Active        *dbo4teamus.SpaceMeetingInfo `json:"active,omitempty" firestore:"active,omitempty"`
}

func (v *RetroSpaceDbo) Validate() error {
	if err := v.RetroSettings.Validate(); err != nil {
		return err
	}
	if v.UpcomingRetro != nil {
		if err := v.UpcomingRetro.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("upcomingRetro", err.Error())
		}
	}
	return nil
}

// ActiveRetro returns info on active retro
func (v *RetroSpaceDbo) ActiveRetro() dbo4teamus.SpaceMeetingInfo {
	if v.Active != nil {
		return *v.Active
	}
	return dbo4teamus.SpaceMeetingInfo{}
}
