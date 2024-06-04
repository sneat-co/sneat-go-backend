package dbo4linkage

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
)

type RelationshipRoleID = string

type RelationshipRole struct {
	//with.CreatedField
}

func (v RelationshipRole) Validate() error {
	return nil
	//return v.CreatedField.Validate()
}

type RelationshipRoles = map[RelationshipRoleID]*RelationshipRole

type RelatedItemKey struct {
	TeamID string `json:"teamID" firestore:"teamID"`
	ItemID string `json:"itemID" firestore:"itemID"`
}

func (v RelatedItemKey) String() string {
	return fmt.Sprintf("%s@%s", v.ItemID, v.TeamID)
}

func (v RelatedItemKey) Validate() error {
	if v.TeamID == "" {
		return validation.NewErrRecordIsMissingRequiredField("teamID")
	}
	if v.ItemID == "" {
		return validation.NewErrRecordIsMissingRequiredField("itemID")
	}
	return nil
}

func HasRelatedItem(relatedItems []*RelatedItem, key RelatedItemKey) bool {
	for _, relatedItem := range relatedItems {
		for _, k := range relatedItem.Keys {
			if k == key {
				return true
			}
		}
	}
	return false
}

func GetRelatedItemByRef(related RelatedByModuleID, itemRef TeamModuleItemRef, createIfMissing bool) *RelatedItem {
	relatedByCollection := related[itemRef.ModuleID]
	if !createIfMissing && len(relatedByCollection) == 0 {
		return nil
	}
	relatedByItem := relatedByCollection[itemRef.Collection]
	if !createIfMissing && len(relatedByItem) == 0 {
		return nil
	}
	for _, ri := range relatedByItem {
		for _, k := range ri.Keys {
			if k.TeamID == itemRef.TeamID && k.ItemID == itemRef.ItemID {
				return ri
			}

		}
	}
	if createIfMissing {
		relatedItem := NewRelatedItem(RelatedItemKey{TeamID: itemRef.TeamID, ItemID: itemRef.ItemID})
		relatedByItem = append(relatedByItem, relatedItem)
		if relatedByCollection == nil {
			relatedByCollection = make(RelatedByCollectionID, 1)
		}
		relatedByCollection[itemRef.Collection] = relatedByItem
		if related == nil {
			related = make(RelatedByModuleID, 1)
		}
		related[itemRef.ModuleID] = relatedByCollection
		return relatedItem
	}
	return nil
}

func GetRelatedItemByKey(relatedItems []*RelatedItem, key RelatedItemKey) *RelatedItem {
	for _, relatedItem := range relatedItems {
		for _, k := range relatedItem.Keys {
			if k == key {
				return relatedItem
			}
		}
	}
	return nil
}

type RelatedItem struct {
	Keys []RelatedItemKey `json:"keys" firestore:"keys"` // TODO: document why we need multiple keys, provide a use case

	Note string `json:"note,omitempty" firestore:"note,omitempty"`

	// RolesOfItem - if related item is a child of the current record, then rolesOfItem = {"child": ...}
	RolesOfItem RelationshipRoles `json:"rolesOfItem,omitempty" firestore:"rolesOfItem,omitempty"`

	// RolesToItem - if related item is a child of the current contact, then rolesToItem = {"parent": ...}
	RolesToItem RelationshipRoles `json:"rolesToItem,omitempty" firestore:"rolesToItem,omitempty"`
}

func (v *RelatedItem) String() string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("RelatedItem{RolesOfItem=%+v, RolesToItem=%+v}", v.RolesOfItem, v.RolesToItem)
}

func NewRelatedItem(key RelatedItemKey) *RelatedItem {
	return &RelatedItem{
		Keys: []RelatedItemKey{key},
	}
}

func (v *RelatedItem) Validate() error {
	if len(v.Keys) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("keys")
	}
	for i, key := range v.Keys {
		if err := key.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("keys[%d]", i), err.Error())
		}
	}
	if err := v.validateRelationships(v.RolesOfItem); err != nil {
		return validation.NewErrBadRecordFieldValue("rolesOfItem", err.Error())
	}
	if err := v.validateRelationships(v.RolesToItem); err != nil {
		return validation.NewErrBadRecordFieldValue("rolesToItem", err.Error())
	}
	return nil
}

func (*RelatedItem) validateRelationships(related RelationshipRoles) error {
	for relationshipID, relationshipDetails := range related {
		if strings.TrimSpace(relationshipID) == "" {
			return validation.NewValidationError("relationship key is empty string")
		}
		if err := relationshipDetails.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(relationshipID, err.Error())
		}
	}
	return nil
}

type RelatedByCollectionID = map[string][]*RelatedItem
type RelatedByModuleID = map[string]RelatedByCollectionID

const relatedField = "related"

var _ Relatable = (*WithRelatedAndIDs)(nil)

func (v *WithRelatedAndIDs) GetRelated() *WithRelatedAndIDs {
	return v
}

type WithRelated struct {
	// Related defines relationships of the current contact to other contacts.
	// Key is team ID.
	Related RelatedByModuleID `json:"related,omitempty" firestore:"related,omitempty"`
}

func (v *WithRelated) Validate() error {
	return v.ValidateRelated(nil)
}

// RemoveRelatedItem removes all relationships to a given item
// TODO(help-wanted): needs 100% code coverage by tests
func (v *WithRelated) RemoveRelatedItem(ref TeamModuleItemRef) (updates []dal.Update) {
	collectionsRelated := v.Related[ref.ModuleID]
	if collectionsRelated == nil {
		return
	}
	relatedItems := collectionsRelated[ref.Collection]
	if len(relatedItems) == 0 {
		return
	}
	var newRelatedItems []*RelatedItem

relatedItemCycle:
	for _, relatedItem := range relatedItems {
		for _, key := range relatedItem.Keys {
			if key.TeamID == ref.TeamID && key.ItemID == ref.ItemID {
				continue relatedItemCycle
			}
		}
		newRelatedItems = append(newRelatedItems, relatedItem)
	}
	field := fmt.Sprintf("%s.%s.%s", relatedField, ref.ModuleID, ref.Collection)
	if len(newRelatedItems) != len(relatedItems) {
		if len(newRelatedItems) == 0 {
			delete(collectionsRelated, ref.Collection)
			if len(collectionsRelated) == 0 {
				delete(v.Related, ref.ModuleID)
				if len(v.Related) == 0 {
					updates = append(updates, dal.Update{
						Field: relatedField,
						Value: dal.DeleteField,
					})
				} else {
					updates = append(updates, dal.Update{
						Field: fmt.Sprintf("%s.%s", relatedField, ref.ModuleID),
						Value: dal.DeleteField,
					})
				}
			} else {
				updates = append(updates, dal.Update{
					Field: field,
					Value: dal.DeleteField,
				})
			}
		} else {
			collectionsRelated[ref.Collection] = newRelatedItems
			updates = append(updates, dal.Update{
				Field: field,
				Value: newRelatedItems,
			})
		}
	}
	return updates
}

func (v *WithRelated) ValidateRelated(validateID func(relatedID string) error) error {
	for moduleID, relatedByCollectionID := range v.Related {
		if moduleID == "" {
			return validation.NewErrBadRecordFieldValue(relatedField, "has empty module ID")
		}
		for collectionID, relatedItems := range relatedByCollectionID {
			if collectionID == "" {
				return validation.NewErrBadRecordFieldValue(
					fmt.Sprintf("%s.%s", relatedField, moduleID),
					"has empty collection ID",
				)
			}
			for i, relatedItem := range relatedItems {
				field := fmt.Sprintf("%s.%s.%s[%d]", relatedField, moduleID, collectionID, i)

				if err := relatedItem.Validate(); err != nil {
					return validation.NewErrBadRecordFieldValue(field, err.Error())
				}
				for _, key := range relatedItem.Keys {
					if validateID != nil {
						relatedID := NewTeamModuleItemRef(key.TeamID, moduleID, collectionID, key.ItemID).ID()
						if err := validateID(relatedID); err != nil {
							return validation.NewErrBadRecordFieldValue(field, err.Error())
						}
					}
				}
			}
		}
	}
	return nil
}

func (v *WithRelated) AddRelationship(itemRef TeamModuleItemRef, rolesCommand RelationshipRolesCommand) (updates []dal.Update, err error) {
	if err := rolesCommand.Validate(); err != nil {
		return nil, err
	}
	if v.Related == nil {
		v.Related = make(RelatedByModuleID, 1)
	}

	if rolesCommand.Add != nil {
		addOppositeRoles := func(roles []RelationshipRoleID, oppositeRoles []RelationshipRoleID) []RelationshipRoleID {
			for _, roleOfItem := range roles {
				if oppositeRole := GetOppositeRole(roleOfItem); oppositeRole != "" && !slice.Contains(rolesCommand.Add.RolesToItem, oppositeRole) {
					oppositeRoles = append(oppositeRoles, oppositeRole)
				}
			}
			return oppositeRoles
		}
		rolesCommand.Add.RolesToItem = addOppositeRoles(rolesCommand.Add.RolesOfItem, rolesCommand.Add.RolesToItem)
		rolesCommand.Add.RolesOfItem = addOppositeRoles(rolesCommand.Add.RolesToItem, rolesCommand.Add.RolesOfItem)
	}

	relatedByCollectionID := v.Related[itemRef.ModuleID]
	if relatedByCollectionID == nil {
		relatedByCollectionID = make(RelatedByCollectionID, 1)
		v.Related[const4contactus.ModuleID] = relatedByCollectionID
	}

	relatedItems := relatedByCollectionID[const4contactus.ContactsCollection]

	relatedItemKey := RelatedItemKey{TeamID: itemRef.TeamID, ItemID: itemRef.ItemID}
	relatedItem := GetRelatedItemByKey(relatedItems, relatedItemKey)
	if relatedItem == nil {
		relatedItem = NewRelatedItem(relatedItemKey)
		relatedItems = append(relatedItems, relatedItem)
		relatedByCollectionID[const4contactus.ContactsCollection] = relatedItems
	}

	addRelationship := func(field string, relationshipIDs []RelationshipRoleID, relationships RelationshipRoles) RelationshipRoles {
		if len(relationshipIDs) == 0 {
			return relationships
		}
		if relationships == nil {
			relationships = make(RelationshipRoles, len(relationshipIDs))
		}
		for _, relationshipID := range relationshipIDs {
			if relationship := relationships[relationshipID]; relationship == nil {
				relationship = &RelationshipRole{
					//CreatedField: with.CreatedField{
					//	Created: with.Created{
					//		By: userID,
					//		At: now.Format(time.RFC3339),
					//	},
					//},
				}
				relationships[relationshipID] = relationship
			}
		}
		return relationships
	}

	if rolesCommand.Add != nil {
		relatedItem.RolesOfItem = addRelationship("rolesOfItem", rolesCommand.Add.RolesOfItem, relatedItem.RolesOfItem)
		relatedItem.RolesToItem = addRelationship("rolesToItem", rolesCommand.Add.RolesToItem, relatedItem.RolesToItem)
	}

	updates = append(updates, dal.Update{
		Field: fmt.Sprintf("related.%s", itemRef.ModuleCollectionPath()),
		Value: relatedItems,
	})

	return updates, nil
}

//func (v *WithRelated) SetRelationshipToItem(
//	userID string,
//	link RelationshipRolesCommand,
//	now time.Time,
//) (updates []dal.Update, err error) {
//	if err = link.Validate(); err != nil {
//		return nil, fmt.Errorf("failed to validate link: %w", err)
//	}
//
//	//var alreadyHasRelatedAs bool
//
//	changed := false
//
//	if v.Related == nil {
//		v.Related = make(RelatedByModuleID, 1)
//	}
//	relatedByCollectionID := v.Related[link.ModuleID]
//	if relatedByCollectionID == nil {
//		relatedByCollectionID = make(RelatedByCollectionID, 1)
//		v.Related[link.ModuleID] = relatedByCollectionID
//	}
//	relatedItems := relatedByCollectionID[const4contactus.ContactsCollection]
//	//if relatedItems == nil {
//	//	relatedItems = make([]*RelatedItem, 0, 1)
//	//	relatedByCollectionID[const4contactus.ContactsCollection] = relatedItems
//	//}
//	relatedItemKey := RelatedItemKey{TeamID: link.TeamID, ItemID: link.ItemID}
//	relatedItem := GetRelatedItemByKey(relatedItems, relatedItemKey)
//	if relatedItem == nil {
//		relatedItem = NewRelatedItem(relatedItemKey)
//		relatedItems = append(relatedItems, relatedItem)
//		relatedByCollectionID[const4contactus.ContactsCollection] = relatedItems
//		changed = true
//	}
//
//	//addIfNeeded := func(f string, itemRelationships RelationshipRoles, linkRelationshipIDs []RelationshipRoleID) {
//	//	field := func() string {
//	//		return fmt.Sprintf("%s.%s.%s", relatedField, link.ID(), f)
//	//	}
//	//	for _, linkRelationshipID := range linkRelationshipIDs {
//	//		itemRelationship := itemRelationships[linkRelationshipID]
//	//		if itemRelationship == nil {
//	//			itemRelationships[linkRelationshipID] = &RelationshipRole{
//	//				CreatedField: with.CreatedField{
//	//					Created: with.Created{
//	//						By: userID,
//	//						At: now.Format(time.DateOnly),
//	//					},
//	//				},
//	//			}
//	//			alreadyHasRelatedAs = true
//	//		}
//	//	}
//	//}
//	//addIfNeeded("rolesOfItem", relatedItem.RolesOfItem, link.RolesOfItem)
//	//addIfNeeded("rolesToItem", relatedItem.RolesToItem, link.RolesToItem)
//
//	var relationshipUpdate []dal.Update
//	if relationshipUpdate, err = v.AddRelationshipAndID(userID, link, now); err != nil {
//		return updates, err
//	}
//	updates = append(updates, relationshipUpdate...)
//
//	return updates, err
//}
