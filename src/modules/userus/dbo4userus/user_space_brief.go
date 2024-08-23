package dbo4userus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/strongo/validation"
)

// UserSpaceBrief hold info on a team in the UserDbo record
type UserSpaceBrief struct {
	dbo4spaceus.SpaceBrief

	// UserContactID is a contact ContactID of a user in the team
	UserContactID string `json:"userContactID" firestore:"userContactID"`

	// UserEntry roles in the team
	Roles []string `json:"roles" firestore:"roles"`

	//MemberType    string   `json:"memberType" firestore:"memberType"` // TODO: document what it is

	// TODO: RetroItems should be moved into members
	//RetroItems dbretro.RetroItemsByType `json:"retroItem,omitempty" firestore:"retroItems,omitempty"`
}

// HasRole checks if a user has a role
func (v UserSpaceBrief) HasRole(role string) bool {
	for _, r := range v.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// Validate validates user record
func (v UserSpaceBrief) Validate() error {
	//if err := models.ValidateTitle(v.Title); err != nil {
	//	return err
	//}
	if v.UserContactID == "" {
		return validation.NewErrRecordIsMissingRequiredField("userContactID")
	}
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	//if v.MemberType == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("memberType")
	//}
	if !core4spaceus.IsValidSpaceType(v.Type) {
		return validation.NewErrBadRecordFieldValue("type", "unknown team type")
	}
	if len(v.Roles) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("roles")
	}
	for i, role := range v.Roles {
		if role == "" {
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("roles[%d]", i))
		}
		if !const4contactus.IsKnownSpaceMemberRole(role, nil) {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("roles[%d]", i), fmt.Sprintf("unknown role (expected one of this role: %+v): %s",
				const4contactus.SpaceMemberWellKnownRoles, role))
		}
	}
	//if len(v.RetroItems) > 0 {
	//	itemIDs := make([]string, len(v.RetroItems)*2) // why 2? why not?
	//	for itemType, items := range v.RetroItems {
	//		if s := strings.TrimSpace(itemType); s == "" {
	//			return validation.NewErrBadRecordFieldValue("retroItems", "retro item with empty item type")
	//		} else if s != itemType || strings.Contains(itemType, " ") {
	//			return validation.NewErrBadRecordFieldValue("retroItems", "spaces in item type")
	//		}
	//		for i, item := range items {
	//			newItemErr := func(message string) error {
	//				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("retroItems[%s][%d]", itemType, i), message)
	//			}
	//			if item == nil {
	//				return newItemErr("nil item")
	//			}
	//			if err := item.Validate(); err != nil {
	//				return newItemErr(err.Error())
	//			}
	//			for _, itemID := range itemIDs {
	//				if itemID == item.ContactID {
	//					return newItemErr("duplicate item ContactID")
	//				}
	//				itemIDs = append(itemIDs, item.ContactID)
	//			}
	//		}
	//	}
	//}
	return nil
}
