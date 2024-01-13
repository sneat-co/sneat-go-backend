package briefs4contactus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
)

type WithContactIDs struct {
	ContactIDs []string `json:"contactIDs,omitempty" firestore:"contactIDs,omitempty"`
}

func (v *WithContactIDs) Validate() error {
	if len(v.ContactIDs) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("contactIDs")
	}
	if v.ContactIDs[0] != "*" {
		return validation.NewErrBadRecordFieldValue("contactIDs[0]", "should be '*'")
	}
	return nil
}

type WithSingleTeamContactIDs struct {
	WithContactIDs
}

func (v *WithSingleTeamContactIDs) AddContactID(contactID string) {
	if len(v.ContactIDs) == 0 {
		v.ContactIDs = make([]string, 1, 2)
		v.ContactIDs[0] = "*"
	}
	v.ContactIDs = append(v.ContactIDs, contactID)
}

func (v *WithSingleTeamContactIDs) HasContactID(contactID string) bool {
	return slice.Contains(v.ContactIDs, contactID)
}

func (v *WithSingleTeamContactIDs) Validate() error {

	return nil
}

// WithMultiTeamContactIDs mixin that adds ContactIDs field
type WithMultiTeamContactIDs struct {
	WithContactIDs
}

func (v *WithMultiTeamContactIDs) AddTeamContactID(teamContactID dbmodels.TeamItemID) {
	v.ContactIDs = append(v.ContactIDs, string(teamContactID))
}

// Validate  returns error if not valid
func (v *WithMultiTeamContactIDs) Validate() error {
	if err := v.WithContactIDs.Validate(); err != nil {
		return err
	}
	for i, id := range v.ContactIDs[1:] {
		switch strings.TrimSpace(string(id)) {
		case string(id): // OK - as expected
		case "":
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contactIDs[%v]", i), "can not be empty string")
		default:
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contactIDs[%v]", i), "leading or trailing whitespaces")
		}
		ids := strings.Split(string(id), dbmodels.TeamItemIDSeparator)
		if len(ids) != 2 {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contactIDs[%v]", i),
				fmt.Sprintf("should be in format '%s', got: %s", dbmodels.NewTeamItemID("teamID", "contactID"), id))
		}
		if ids[0] == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contactIDs[%v]", i), "teamID can not be empty string")
		}
		if ids[1] == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contactIDs[%v]", i), "contactID can not be empty string")
		}
	}
	return nil
}

// HasTeamContactID check if a record has a specific contactBrief ID
func (v *WithMultiTeamContactIDs) HasTeamContactID(teamItemID dbmodels.TeamItemID) bool {
	return slice.Contains(v.ContactIDs, string(teamItemID))
}
