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

type TeamModuleDocRef struct { // TODO: Move to sneat-go-core or document why not
	TeamID     string `json:"teamID" firestore:"teamID"`
	ModuleID   string `json:"moduleID" firestore:"moduleID"`
	Collection string `json:"collection" firestore:"collection"`
	ItemID     string `json:"itemID" firestore:"itemID"`
}

func NewTeamModuleDocRef(teamID, moduleID, collection, itemID string) TeamModuleDocRef {
	return TeamModuleDocRef{
		TeamID:     teamID,
		ModuleID:   moduleID,
		Collection: collection,
		ItemID:     itemID,
	}
}

func NewTeamModuleDocRefFromString(id string) TeamModuleDocRef {
	ids := strings.Split(id, ".")
	return TeamModuleDocRef{
		TeamID:     ids[0],
		ModuleID:   ids[1],
		Collection: ids[2],
		ItemID:     ids[3],
	}
}

func (v TeamModuleDocRef) ID() string {
	return fmt.Sprintf("%s.%s.%s.%s", v.TeamID, v.ModuleID, v.Collection, v.ItemID)
}

func (v TeamModuleDocRef) Validate() error {
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
	TeamModuleDocRef
	//
	RelatedAs []RelationshipID `json:"relatedAs,omitempty" firestore:"relatedAs,omitempty"`
	RelatesAs []RelationshipID `json:"relatesAs,omitempty" firestore:"relatesAs,omitempty"`
}

func (v Link) Validate() error {
	if err := v.TeamModuleDocRef.Validate(); err != nil {
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
	if err := valRelationIDs("relatedAs", v.RelatedAs); err != nil {
		return err
	}
	if err := valRelationIDs("relatesAs", v.RelatesAs); err != nil {
		return err
	}
	return nil
}
