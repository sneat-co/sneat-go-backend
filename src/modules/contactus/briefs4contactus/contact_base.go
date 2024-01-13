package briefs4contactus

import (
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

// ContactBase is used in models4contactus.ContactDto and in requests to create a contactBrief
type ContactBase struct {
	ContactBrief
	dbmodels.WithUpdatedAndVersion

	// Status belong to ContactBase and is not part of ContactBrief as we keep in briefs only active contacts
	Status dbmodels.Status `json:"status" firestore:"status"` // active, archived

	WithGroupIDs

	Address   *dbmodels.Address `json:"address,omitempty" firestore:"address,omitempty"`
	VATNumber string            `json:"vatNumber,omitempty" firestore:"vatNumber,omitempty"`

	// Dob is Date of birth
	DoB string `json:"dob,omitempty" firestore:"dob,omitempty"`

	Emails []dbmodels.PersonEmail `json:"emails,omitempty" firestore:"emails,omitempty"`
	Phones []dbmodels.PersonPhone `json:"phones,omitempty" firestore:"phones,omitempty"`

	Timezone *dbmodels.Timezone `json:"timezone,omitempty" firestore:"timezone,omitempty"`
}

func (v *ContactBase) Equal(v2 *ContactBase) bool {
	return v.ContactBrief.Equal(&v2.ContactBrief) && v.VATNumber == v2.VATNumber
}

// Validate returns error if not valid
func (v *ContactBase) Validate() error {
	var errs []error
	if err := v.ContactBrief.Validate(); err != nil {
		errs = append(errs, err)
	}
	switch v.Status {
	case ContactStatusActive, ContactStatusArchived, ContactStatusDeleted: // OK
	case "":
		errs = append(errs, validation.NewErrRequestIsMissingRequiredField("status"))
	default:
		errs = append(errs, validation.NewErrBadRecordFieldValue("status", "unknown value: "+v.Status))
	}
	if v.Type == ContactTypeCompany {
		if v.CountryID == "" {
			errs = append(errs, validation.NewErrBadRecordFieldValue("countryID", "missing required field for a contactBrief of type=company"))
		}
	}
	if strings.TrimSpace(v.Title) == "" && v.Names == nil {
		errs = append(errs, validation.NewErrRecordIsMissingRequiredField("name|title"))
	}
	if v.Names != nil {
		if err := v.Names.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	switch v.Type {
	case ContactTypePerson, ContactTypeAnimal:
		if v.VATNumber != "" {
			return validation.NewErrBadRecordFieldValue("vatNumber", "should be empty for a contactBrief of type=person")
		}
		if err := dbmodels.ValidateGender(v.Gender, true); err != nil {
			errs = append(errs, err)
		}
		if v.AgeGroup != "" || v.Type != ContactTypeAnimal {
			if err := dbmodels.ValidateAgeGroup(v.AgeGroup, true); err != nil {
				errs = append(errs, err)
			}
		}
	case ContactTypeCompany:
		if v.Gender != "" {
			errs = append(errs, validation.NewErrBadRecordFieldValue("gender", "expected to be empty for contactBrief type=company, got: "+v.Gender))
		}
		if err := dbmodels.ValidateGender(v.Gender, false); err != nil {
			return err
		}
		if v.AgeGroup != "" {
			errs = append(errs, validation.NewErrBadRecordFieldValue("ageGroup", "expected to be empty for contactBrief type=company, got: "+v.Gender))
		}
		if err := dbmodels.ValidateAgeGroup(v.AgeGroup, false); err != nil {
			errs = append(errs, err)
		}
	}

	hasPrimaryEmail := false
	for i, email := range v.Emails {
		if err := email.Validate(); err != nil {
			errs = append(errs, validation.NewErrBadRecordFieldValue(fmt.Sprintf("emails[%v]", i), err.Error()))
		}
		if email.Type == "primary" {
			if hasPrimaryEmail {
				errs = append(errs, validation.NewErrBadRecordFieldValue("emails", "only 1 email can have type=primary"))
			}
			hasPrimaryEmail = true
		}
	}
	for i, phone := range v.Phones {
		if err := phone.Validate(); err != nil {
			errs = append(errs, validation.NewErrBadRecordFieldValue(fmt.Sprintf("phones[%v]", i), err.Error()))
		}
	}
	if l := len(errs); l == 1 {
		return validation.NewErrBadRecordFieldValue("ContactBase", errs[0].Error())
	} else if l > 0 {
		return validation.NewErrBadRecordFieldValue("ContactBase", fmt.Errorf("%d errors:\n%w", l, errors.Join(errs...)).Error())
	}
	if err := v.WithGroupIDs.Validate(); err != nil {
		return err
	}
	if err := v.Timezone.Validate(); err != nil {
		return fmt.Errorf("invalid 'timezone' field: %w", err)
	}
	return nil
}
