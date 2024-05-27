package briefs4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/models/dbprofile"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

// ContactBrief needed as ContactBase is used in models4contactus.ContactDto and in dto4contactus.CreatePersonRequest
// Status is not part of ContactBrief as we keep in briefs only active contacts
type ContactBrief struct {
	dbmodels.WithUserID
	dbmodels.WithOptionalRelatedAs // This is used in `Related` field of `ContactDbo`
	with.OptionalCountryID
	with.RolesField

	Type       ContactType        `json:"type" firestore:"type"` // "person", "company", "location"
	Gender     dbmodels.Gender    `json:"gender,omitempty" firestore:"gender,omitempty"`
	Names      *person.NameFields `json:"names,omitempty" firestore:"names,omitempty"`
	Title      string             `json:"title,omitempty" firestore:"title,omitempty"`
	ShortTitle string             `json:"shortTitle,omitempty" firestore:"shortTitle,omitempty"` // Not supposed to be used in models4contactus.ContactDto
	ParentID   string             `json:"parentID" firestore:"parentID"`                         // Intentionally not adding `omitempty` so we can search root contacts only

	// Number of active invites to join a team
	InvitesCount int `json:"activeInvitesCount,omitempty" firestore:"activeInvitesCount,omitempty"`

	// AgeGroup is deprecated?
	AgeGroup string `json:"ageGroup,omitempty" firestore:"ageGroup,omitempty"` // TODO: Add validation
	PetKind  string `json:"species,omitempty" firestore:"species,omitempty"`

	// Avatar holds a photo of a member
	Avatar *dbprofile.Avatar `json:"avatar,omitempty" firestore:"avatar,omitempty"`
}

func (v *ContactBrief) SetName(field, value string) {
	if v.Names == nil {
		v.Names = &person.NameFields{}
	}
	switch field {
	case "first":
		v.Names.FirstName = value
	case "last":
		v.Names.LastName = value
	case "middle":
		v.Names.MiddleName = value
	case "full":
		v.Names.FullName = value
	case "nick":
		v.Names.NickName = value
	default:
		panic("unsupported field: " + field)
	}
}

func (v *ContactBrief) IsTeamMember() bool {
	return v.HasRole(const4contactus.TeamMemberRoleMember)
}

// GetUserID returns UserID field value
func (v *ContactBrief) GetUserID() string {
	return v.UserID
}

// Equal returns true if 2 instances are equal
func (v *ContactBrief) Equal(v2 *ContactBrief) bool {
	return v.Type == v2.Type &&
		v.WithUserID == v2.WithUserID &&
		v.Gender == v2.Gender &&
		v.OptionalCountryID == v2.OptionalCountryID &&
		v.Names.Equal(v2.Names) &&
		v.WithOptionalRelatedAs.Equal(v2.WithOptionalRelatedAs) &&
		v.Avatar.Equal(v2.Avatar)
}

// Validate returns error if not valid
func (v *ContactBrief) Validate() error {
	if err := ValidateContactType(v.Type); err != nil {
		return err
	}
	if err := dbmodels.ValidateGender(v.Gender, false); err != nil {
		return err
	}
	if strings.TrimSpace(v.Title) == "" && v.Names == nil {
		return validation.NewErrRecordIsMissingRequiredField("name|title")
	} else if err := v.Names.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("name", err.Error())
	}
	if v.UserID != "" {
		if !core.IsAlphanumericOrUnderscore(v.UserID) {
			return validation.NewErrBadRecordFieldValue("userID", "is not alphanumeric: "+v.UserID)
		}
	}
	switch v.Type {
	case ContactTypeLocation:
		if v.ParentID == "" {
			return validation.NewErrRecordIsMissingRequiredField("parentID")
		}
	}
	if err := v.OptionalCountryID.Validate(); err != nil {
		return err
	}
	if err := v.RolesField.Validate(); err != nil {
		return err
	}
	if err := v.WithUserID.Validate(); err != nil {
		return err
	}
	if v.PetKind != "" {
		if !const4contactus.IsKnownPetPetKind(v.PetKind) {
			return validation.NewErrBadRecordFieldValue("species", "unknown value: "+v.PetKind)
		}
	}
	return nil
}

// GetTitle return full name of a person
func (v *ContactBrief) GetTitle() string {
	if v.Title != "" {
		return v.Title
	}
	return v.Names.GetFullName()
}

func (v *ContactBrief) DetermineShortTitle(title string, contacts map[string]*ContactBrief) string {
	if v.Names.FirstName != "" && IsUniqueShortTitle(v.Names.FirstName, contacts, const4contactus.TeamMemberRoleMember) {
		v.ShortTitle = v.Names.FirstName
	} else if v.Names.NickName != "" && IsUniqueShortTitle(v.Names.FirstName, contacts, const4contactus.TeamMemberRoleMember) {
		return v.Names.NickName
	} else if v.Names.FullName != "" {
		return getShortTitle(v.Names.FullName, contacts)
	} else if title != "" {
		return getShortTitle(title, contacts)
	}
	return ""
}

func getShortTitle(title string, members map[string]*ContactBrief) string {
	shortNames := GetShortNames(title)
	for _, short := range shortNames {
		isUnique := true
		for _, m := range members {
			if m.ShortTitle == short.Name {
				isUnique = false
				break
			}
		}
		if isUnique {
			return short.Name
		}
	}
	return ""
}

type ShortName struct {
	Name string `json:"name" firestore:"name"`
	Type string `json:"type" firestore:"type"`
}

// GetShortNames returns short names from a title
func GetShortNames(title string) (shortNames []ShortName) {
	title = CleanTitle(title)
	names := strings.Split(title, " ")
	shortNames = make([]ShortName, 0, len(names))
NAMES:
	for _, s := range names {
		name := strings.TrimSpace(s)
		if name == "" {
			continue
		}
		for _, sn := range shortNames {
			if sn.Name == name {
				continue NAMES
			}
		}
		shortNames = append(shortNames, ShortName{
			Name: name,
			Type: "unknown",
		})
	}
	return shortNames
}

// CleanTitle cleans title from spaces
func CleanTitle(title string) string {
	title = strings.TrimSpace(title)
	for strings.Contains(title, "  ") {
		title = strings.Replace(title, "  ", " ", -1)
	}
	return title
}
