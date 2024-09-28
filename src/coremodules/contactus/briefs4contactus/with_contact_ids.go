package briefs4contactus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"slices"
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

type WithSingleSpaceContactIDs struct {
	WithContactIDs
}

func (v *WithSingleSpaceContactIDs) AddContactID(contactID string) {
	if len(v.ContactIDs) == 0 {
		v.ContactIDs = make([]string, 1, 2)
		v.ContactIDs[0] = "*"
	}
	v.ContactIDs = append(v.ContactIDs, contactID)
}

func (v *WithSingleSpaceContactIDs) HasContactID(contactID string) bool {
	return slices.Contains(v.ContactIDs, contactID)
}

func (v *WithSingleSpaceContactIDs) Validate() error {

	return nil
}

// WithMultiSpaceContactIDs mixin that adds ContactIDs field
type WithMultiSpaceContactIDs struct {
	WithContactIDs
}

func (v *WithMultiSpaceContactIDs) AddSpaceContactID(teamContactID dbmodels.SpaceItemID) {
	v.ContactIDs = append(v.ContactIDs, string(teamContactID))
}

// Validate  returns error if not valid
func (v *WithMultiSpaceContactIDs) Validate() error {
	if err := v.WithContactIDs.Validate(); err != nil {
		return err
	}
	for i, id := range v.ContactIDs[1:] {
		switch strings.TrimSpace(string(id)) {
		case string(id): // OK - as expected
		case "":
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contactIDs[%d]", i), "can not be empty string")
		default:
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contactIDs[%d]", i), "leading or trailing whitespaces")
		}
		ids := strings.Split(string(id), dbmodels.SpaceItemIDSeparator)
		if len(ids) != 2 {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contactIDs[%d]", i),
				fmt.Sprintf("should be in format '%s', got: %s", dbmodels.NewSpaceItemID("spaceID", "contactID"), id))
		}
		if ids[0] == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contactIDs[%d]", i), "spaceID can not be empty string")
		}
		if ids[1] == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("contactIDs[%d]", i), "contactID can not be empty string")
		}
	}
	return nil
}

// HasSpaceContactID check if a record has a specific contactBrief ContactID
func (v *WithMultiSpaceContactIDs) HasSpaceContactID(teamItemID dbmodels.SpaceItemID) bool {
	return slices.Contains(v.ContactIDs, string(teamItemID))
}
