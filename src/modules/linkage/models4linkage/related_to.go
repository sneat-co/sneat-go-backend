package models4linkage

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/slice"
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

type RolesCommand struct {
	RolesOfItem []RelationshipRoleID `json:"rolesOfItem,omitempty" firestore:"rolesOfItem,omitempty"`
	RolesToItem []RelationshipRoleID `json:"rolesToItem,omitempty" firestore:"rolesToItem,omitempty"`
}

type Link struct {
	TeamModuleItemRef
	Add    *RolesCommand `json:"add,omitempty" firestore:"add,omitempty"`
	Remove *RolesCommand `json:"remove,omitempty" firestore:"remove,omitempty"`
}

func (v Link) Validate() error {
	if err := v.TeamModuleItemRef.Validate(); err != nil {
		return err
	}
	if err := v.Add.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("add", err.Error())
	}
	if err := v.Remove.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("remove", err.Error())
	}
	return nil
}
func (v *RolesCommand) Validate() error {
	if v == nil {
		return nil
	}
	validateRelationIDs := func(field string, relations []string) error {
		for i, s := range relations {
			if strings.TrimSpace(s) != s {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("%s[%d]", field, i),
					"must not have leading or trailing spaces")
			}
			if slice.Contains(relations[:i], s) {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("%s[%d]", field, i),
					"duplicate relationship role value: "+s)
			}
		}
		return nil
	}
	if v.RolesToItem == nil && v.RolesOfItem == nil {
		return validation.NewErrRecordIsMissingRequiredField("rolesOfItem|rolesToItem")
	}
	if v.RolesToItem != nil {
		if err := validateRelationIDs("rolesOfItem", v.RolesOfItem); err != nil {
			return err
		}
	}
	if v.RolesToItem != nil {
		if err := validateRelationIDs("rolesToItem", v.RolesToItem); err != nil {
			return err
		}
	}
	return nil
}
