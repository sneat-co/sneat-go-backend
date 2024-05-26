package models4linkage

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strings"
)

type ShortTeamModuleDocRef struct {
	ID     string `json:"id" firestore:"id"`
	TeamID string `json:"teamID,omitempty" firestore:"teamID,omitempty"`
}

func (v *ShortTeamModuleDocRef) Validate() error {
	// TeamID can be empty for global collections like Happening
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("itemID")
	} else if err := validate.RecordID(v.ID); err != nil {
		return validation.NewErrBadRecordFieldValue("itemID", err.Error())
	}
	return nil
}

type TeamModuleItemRef struct { // TODO: Move to sneat-go-core or document why not
	TeamID     string `json:"teamID" firestore:"teamID"`
	ModuleID   string `json:"moduleID" firestore:"moduleID"`
	Collection string `json:"collection" firestore:"collection"`
	ItemID     string `json:"itemID" firestore:"itemID"`
}

func NewTeamModuleDocRef(teamID, moduleID, collection, itemID string) TeamModuleItemRef {
	return TeamModuleItemRef{
		TeamID:     teamID,
		ModuleID:   moduleID,
		Collection: collection,
		ItemID:     itemID,
	}
}

func NewTeamModuleDocRefFromString(id string) TeamModuleItemRef {
	ids := strings.Split(id, ".")
	if len(ids) != 4 {
		panic(fmt.Sprintf("invalid ID: '%s'", id))
	}
	return TeamModuleItemRef{
		ModuleID:   ids[0],
		Collection: ids[1],
		TeamID:     ids[2],
		ItemID:     ids[3],
	}
}

func (v TeamModuleItemRef) ID() string {
	return fmt.Sprintf("%s.%s.%s", v.ModuleCollectionPath(), v.TeamID, v.ItemID)
}

func (v TeamModuleItemRef) ModuleCollectionPath() string {
	return fmt.Sprintf("%s.%s", v.ModuleID, v.Collection)
}

func (v TeamModuleItemRef) Validate() error {
	// TeamID can be empty for global collections like Happening
	if v.ModuleID == "" {
		return validation.NewErrRecordIsMissingRequiredField("moduleID")
	}
	if v.Collection == "" {
		return validation.NewErrRecordIsMissingRequiredField("collection")
	}
	if v.ItemID == "" {
		return validation.NewErrRecordIsMissingRequiredField("itemID")
	} else if err := validate.RecordID(v.ItemID); err != nil {
		return validation.NewErrBadRecordFieldValue("itemID", err.Error())
	}
	return nil
}

type Link struct {
	TeamModuleItemRef
	//
	RolesOfItem []RelationshipRoleID `json:"rolesOfItem,omitempty" firestore:"rolesOfItem,omitempty"`
	RolesToItem []RelationshipRoleID `json:"rolesToItem,omitempty" firestore:"rolesToItem,omitempty"`
}

func (v Link) Validate() error {
	if err := v.TeamModuleItemRef.Validate(); err != nil {
		return err
	}
	valRelationIDs := func(field string, relations []string) error {
		for i, s := range relations {
			if strings.TrimSpace(s) != s {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("%s[%d]", field, i),
					"must not have leading or trailing spaces")
			}
		}
		return nil
	}
	if err := valRelationIDs("rolesOfItem", v.RolesOfItem); err != nil {
		return err
	}
	if err := valRelationIDs("rolesToItem", v.RolesToItem); err != nil {
		return err
	}
	return nil
}
