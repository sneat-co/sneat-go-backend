package dbo4userus

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/core/coremodels"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/const4userus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"net/mail"
	"slices"
	"strings"
	"time"
)

type WithUserIDs struct {
	UserIDs map[string]string `json:"userIDs,omitempty" firestore:"userIDs,omitempty"`
}

func (v *WithUserIDs) SetUserID(spaceID string, userID string) {
	if v.UserIDs == nil {
		v.UserIDs = map[string]string{spaceID: userID}
	} else {
		v.UserIDs[spaceID] = userID
	}
}

var _ botsfwmodels.AppUserData = (*UserDbo)(nil)
var _ botsfwmodels.AppUserData = (*userBotsFwAdapter)(nil)

// UserDbo is a record that holds information about user
type UserDbo struct {
	briefs4contactus.ContactBase
	with.CreatedFields
	dbmodels.WithPreferredLocale
	dbmodels.WithPrimaryCurrency
	dbmodels.WithLastCurrencies
	botsfwmodels.WithBotUserIDs
	coremodels.SmsStats

	appuser.AccountsOfUser
	appuser.WithLastLogin

	dbo4linkage.WithRelatedAndIDs

	InvitedByUserID string `firestore:"invitedByUserID,omitempty" ` // TODO: Prevent circular references! see users 6032980589936640 & 5998019824582656

	IsAnonymous bool `json:"isAnonymous" firestore:"isAnonymous"` // Intentionally do not omitempty
	//Title string `json:"title,omitempty" firestore:"title,omitempty"`

	Timezone *dbmodels.Timezone `json:"timezone,omitempty" firestore:"timezone,omitempty"`

	Defaults *UserDefaults `json:"defaults,omitempty" firestore:"defaults,omitempty"`

	Email         string `json:"email,omitempty"  firestore:"email,omitempty"`
	EmailVerified bool   `json:"emailVerified"  firestore:"emailVerified"`

	// List of teams a user belongs to
	Spaces   map[string]*UserSpaceBrief `json:"spaces,omitempty"   firestore:"spaces,omitempty"`
	SpaceIDs []string                   `json:"spaceIDs,omitempty" firestore:"spaceIDs,omitempty"`

	Created dbmodels.CreatedInfo `json:"created" firestore:"created"`

	//models4debtus.WithGroups

	// TODO: Should this be moved to company members?
	//models.DatatugUser

	ReferredBy string `firestore:"referredBy,omitempty"`

	LastFeedbackAt   time.Time `firestore:"lastFeedbackAt,omitempty"`
	LastFeedbackRate string    `firestore:"lastFeedbackRate,omitempty"`
}

func (v *UserDbo) GetFullName() string {
	return v.Names.GetFullName()
}

// SetSpaceBrief sets team brief and adds teamID to the list of team IDs if needed
func (v *UserDbo) SetSpaceBrief(spaceID string, brief *UserSpaceBrief) (updates []dal.Update) {
	if spaceID == "" {
		panic("spaceID is empty string")
	}
	if brief == nil {
		panic("brief is nil")
	}
	if v.Spaces == nil {
		v.Spaces = make(map[string]*UserSpaceBrief, 1)
	}
	v.Spaces[spaceID] = brief
	updates = append(updates, dal.Update{Field: "spaces." + spaceID, Value: brief})
	if !slices.Contains(v.SpaceIDs, spaceID) {
		v.SpaceIDs = append(v.SpaceIDs, spaceID)
		updates = append(updates, dal.Update{Field: "spaceIDs", Value: v.SpaceIDs})
	}
	return
}

func (v *UserDbo) GetFamilySpaceID() string {
	id, _ := v.GetFirstSpaceBriefBySpaceType(core4spaceus.SpaceTypeFamily)
	return id
}

// GetSpaceBriefsByType returns the all spaces matching a specific type
func (v *UserDbo) GetSpaceBriefsByType(t core4spaceus.SpaceType) (spaces map[string]*UserSpaceBrief) {
	for id, brief := range v.Spaces {
		if brief.Type == t {
			if spaces == nil {
				spaces = make(map[string]*UserSpaceBrief)
			}
			spaces[id] = brief
		}
	}
	return
}

func (v *UserDbo) GetFirstSpaceBriefBySpaceType(spaceType core4spaceus.SpaceType) (spaceID string, spaceBrief *UserSpaceBrief) {
	for id, space := range v.Spaces {
		if space.Type == spaceType {
			return id, space
		}
	}
	return "", nil
}

// Validate validates user record
func (v *UserDbo) Validate() error {
	if err := v.ContactBase.Validate(); err != nil {
		return err
	}
	if err := v.SmsStats.Validate(); err != nil {
		return err
	}
	//if v.Avatar != nil {
	//	if err := v.Avatar.Validate(); err != nil {
	//		return validation.NewErrBadRecordFieldValue("avatar", err.Error())
	//	}
	//}
	//if v.Title != "" {
	//	if err := v.Names.Validate(); err != nil {
	//		return err
	//	}
	//}
	if err := v.validateEmails(); err != nil {
		return err
	}
	if err := v.validateSpaces(); err != nil {
		return err
	}
	if err := dbmodels.ValidateGender(v.Gender, true); err != nil {
		return err
	}
	//if v.Datatug != nil {
	//	if err := v.Datatug.Validate(); err != nil {
	//		return err
	//	}
	//}
	if err := v.Created.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("created", err.Error())
	}
	if err := v.WithRelatedAndIDs.Validate(); err != nil {
		return err
	}
	return nil
}

func (v *UserDbo) validateEmails() error {
	if strings.TrimSpace(v.Email) != v.Email {
		return validation.NewErrBadRecordFieldValue("email", "contains leading or closing spaces")
	}
	if strings.Contains(v.Email, " ") {
		return validation.NewErrBadRecordFieldValue("email", "contains space")
	}
	if v.Email != "" {
		if _, err := mail.ParseAddress(v.Email); err != nil {
			return validation.NewErrBadRecordFieldValue("email", err.Error())
		}
		if len(v.Emails) == 0 {
			return validation.NewErrBadRecordFieldValue("emails", "user record has 'email' value but 'emails' are empty")
		}
	}
	primaryEmailInEmails := false
	for i, email := range v.Emails {
		if err := email.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("emails[%d]", i), err.Error())
		}
		if email.Address == v.Email {
			primaryEmailInEmails = true
		}
	}
	if v.Email != "" && !primaryEmailInEmails {
		return validation.NewErrBadRecordFieldValue("emails", "user's primary email is not in 'emails' field")
	}
	return nil
}

func (v *UserDbo) validateSpaces() error {
	if len(v.Spaces) != len(v.SpaceIDs) {
		return validation.NewErrBadRecordFieldValue("spaceIDs",
			fmt.Sprintf("len(v.Spaces) != len(v.SpaceIDs): %d != %d", len(v.Spaces), len(v.SpaceIDs)))
	}
	if len(v.Spaces) > 0 {
		spaceIDs := make([]string, 0, len(v.Spaces))
		spaceTitles := make([]string, 0, len(v.Spaces))
		for spaceID, space := range v.Spaces {
			if spaceID == "" {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("spaces['%s']", spaceID), "holds empty id")
			}
			if !slices.Contains(v.SpaceIDs, spaceID) {
				return validation.NewErrBadRecordFieldValue("spaceIDs", "missing team ContactID: "+spaceID)
			}
			if err := space.Validate(); err != nil {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("spaces[%s]{title=%s}", spaceID, space.Title), err.Error())
			}
			if space.Title != "" {
				if i := slices.Index(spaceTitles, space.Title); i >= 0 {
					return validation.NewErrBadRecordFieldValue("spaces",
						fmt.Sprintf("at least 2 spaces (%s & %s) with same title: '%s'", spaceID, spaceIDs[i], space.Title))
				}
			}
			spaceIDs = append(spaceIDs, spaceID)
			spaceTitles = append(spaceTitles, space.Title)
		}
	}
	return nil
}

// GetUserSpaceInfoByID returns team info specific to the user by team ContactID
func (v *UserDbo) GetUserSpaceInfoByID(spaceID string) *UserSpaceBrief {
	return v.Spaces[spaceID]
}

func (v *UserDbo) SetBotUserID(platform const4userus.AuthProviderCode, botID, botUserID string) {
	v.AddAccount(appuser.AccountKey{
		Provider: platform,
		App:      botID,
		ID:       botUserID,
	})
}
