package dbo4linkage

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
)

type ShortSpaceModuleDocRef struct {
	ID      string `json:"id" firestore:"id"`
	SpaceID string `json:"spaceID,omitempty" firestore:"spaceID,omitempty"`
}

func (v *ShortSpaceModuleDocRef) Validate() error {
	// SpaceID can be empty for global collections like Happening
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("itemID")
	} else if err := validate.RecordID(v.ID); err != nil {
		return validation.NewErrBadRecordFieldValue("itemID", err.Error())
	}
	return nil
}

type SpaceModuleItemRef struct { // TODO: Move to sneat-go-core or document why not
	SpaceID    string `json:"spaceID" firestore:"spaceID"`
	ModuleID   string `json:"moduleID" firestore:"moduleID"`
	Collection string `json:"collection" firestore:"collection"`
	ItemID     string `json:"itemID" firestore:"itemID"`
}

func NewSpaceModuleItemRef(teamID, moduleID, collection, itemID string) SpaceModuleItemRef {
	return SpaceModuleItemRef{
		SpaceID:    teamID,
		ModuleID:   moduleID,
		Collection: collection,
		ItemID:     itemID,
	}
}

func NewSpaceModuleItemRefFromString(id string) SpaceModuleItemRef {
	ids := strings.Split(id, ".")
	if len(ids) != 4 {
		panic(fmt.Sprintf("invalid ID: '%s'", id))
	}
	return SpaceModuleItemRef{
		ModuleID:   ids[0],
		Collection: ids[1],
		SpaceID:    ids[2],
		ItemID:     ids[3],
	}
}

func (v SpaceModuleItemRef) ID() string {
	return fmt.Sprintf("%s.%s.%s", v.ModuleCollectionPath(), v.SpaceID, v.ItemID)
}

func (v SpaceModuleItemRef) ModuleCollectionPath() string {
	return fmt.Sprintf("%s.%s", v.ModuleID, v.Collection)
}

func (v SpaceModuleItemRef) Validate() error {
	// SpaceID can be empty for global collections like Happening
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

type RelationshipRolesCommand struct {
	//SpaceModuleItemRef
	Add    *RolesCommand `json:"add,omitempty" firestore:"add,omitempty"`
	Remove *RolesCommand `json:"remove,omitempty" firestore:"remove,omitempty"`
}

func (v RelationshipRolesCommand) Validate() error {
	//if err := v.SpaceModuleItemRef.Validate(); err != nil {
	//	return err
	//}
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
