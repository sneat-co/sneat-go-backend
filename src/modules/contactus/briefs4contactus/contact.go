package briefs4contactus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"slices"
)

// WithGroupIDs is a mixin that adds groupIDs field
type WithGroupIDs struct {
	GroupIDs []string `json:"groupIDs,omitempty" firestore:"groupIDs,omitempty"`
}

// Validate returns error if not valid
func (v WithGroupIDs) Validate() error {
	if err := with.ValidateSetSliceField("groupIDs", v.GroupIDs, true); err != nil {
		return err
	}
	return nil
}

// ContactType is a type of contactBrief that can be: person, company, location, vehicle
type ContactType string

const (
	ContactTypePerson   = "person"
	ContactTypeCompany  = "company"
	ContactTypeLocation = "location"
	//ContactTypeVehicle  = "vehicle"
	ContactTypeAnimal = "animal"
	//ContactTypePet      = "pet"
)

var ContactTypes = []ContactType{
	ContactTypePerson,
	ContactTypeCompany,
	ContactTypeLocation,
	//ContactTypeVehicle,
	ContactTypeAnimal,
	//ContactTypePet,
}

// ValidateContactIDRecordField validates contactBrief ContactID record field
func ValidateContactIDRecordField(name, value string, isRequired bool) error {
	if !isRequired && value == "" {
		return nil
	}
	if err := validate.RecordID(value); err != nil {
		return validation.NewErrBadRecordFieldValue(name, err.Error())
	}
	return nil
}

// ValidateContactType returns error if invalid value
func ValidateContactType(v ContactType) error {
	if v == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if !slices.Contains(ContactTypes, v) {
		return validation.NewErrBadRecordFieldValue("type",
			fmt.Sprintf("unknown value: [%s]", v))
	}
	return nil
}

const (
	ContactStatusActive   = dbmodels.StatusActive
	ContactStatusArchived = dbmodels.StatusArchived
	ContactStatusDeleted  = dbmodels.StatusDeleted
)
