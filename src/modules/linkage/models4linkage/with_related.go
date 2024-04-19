package models4linkage

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
	"time"
)

type RelationshipID = string

type Relationship struct {
	with.CreatedField
}

func (v Relationship) Validate() error {
	return v.CreatedField.Validate()
}

type Relationships = map[RelationshipID]*Relationship

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
	Keys []RelatedItemKey `json:"keys" firestore:"keys"`

	// RelatedAs - if related contact is a child of the current contact, then relatedAs = {"child": ...}
	RelatedAs Relationships `json:"relatedAs,omitempty" firestore:"relatedAs,omitempty"`

	// RelatesAs - if related contact is a child of the current contact, then relatesAs = {"parent": ...}
	RelatesAs Relationships `json:"relatesAs,omitempty" firestore:"relatesAs,omitempty"`
}

func (v *RelatedItem) String() string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("RelatedItem{RelatedAs=%+v, RelatesAs=%+v}", v.RelatedAs, v.RelatesAs)
}

func NewRelatedItem(key RelatedItemKey) *RelatedItem {
	return &RelatedItem{
		Keys:      []RelatedItemKey{key},
		RelatedAs: make(Relationships, 1),
		RelatesAs: make(Relationships, 1),
	}
}

func (v *RelatedItem) Validate() error {
	if len(v.Keys) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("keys")
	}
	for i, key := range v.Keys {
		if err := key.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("keys[%v]", i), err.Error())
		}
	}
	if err := v.validateRelationships(v.RelatedAs); err != nil {
		return validation.NewErrBadRecordFieldValue("relatedAs", err.Error())
	}
	if err := v.validateRelationships(v.RelatesAs); err != nil {
		return validation.NewErrBadRecordFieldValue("relatesAs", err.Error())
	}
	return nil
}

func (*RelatedItem) validateRelationships(related Relationships) error {
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

// RemoveRelationshipToContact removes all relationships to a given contact
func (v *WithRelated) RemoveRelationshipToContact(ref TeamModuleDocRef) (updates []dal.Update) {
	id := ref.ID()
	if _, ok := v.Related[id]; ok {
		delete(v.Related, id)
		updates = append(updates, dal.Update{
			Field: relatedField + "." + id,
			Value: dal.DeleteField,
		})
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
						if err := validateID(key.String()); err != nil {
							return validation.NewErrBadRecordFieldValue(field, err.Error())
						}
					}
				}
			}
		}
	}
	return nil
}

func (v *WithRelated) SetRelationshipToItem(
	userID string,
	link Link,
	now time.Time,
) (updates []dal.Update, err error) {
	if err = link.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate link: %w", err)
	}

	var alreadyHasRelatedAs bool

	if v.Related == nil {
		v.Related = make(RelatedByModuleID, 1)
	}
	relatedByCollectionID := v.Related[link.ModuleID]
	if relatedByCollectionID == nil {
		relatedByCollectionID = make(RelatedByCollectionID, 1)
		v.Related[link.ModuleID] = relatedByCollectionID
	}
	relatedItems := relatedByCollectionID[const4contactus.ContactsCollection]
	if relatedItems == nil {
		relatedItems = make([]*RelatedItem, 0, 1)
		relatedByCollectionID[const4contactus.ContactsCollection] = relatedItems
	}
	relatedItemKey := RelatedItemKey{TeamID: link.TeamID, ItemID: link.ItemID}
	relatedItem := GetRelatedItemByKey(relatedItems, relatedItemKey)
	if relatedItem == nil {
		relatedItem = NewRelatedItem(relatedItemKey)
		relatedItems = append(relatedItems, relatedItem)
		relatedByCollectionID[const4contactus.ContactsCollection] = relatedItems
	}

	addIfNeeded := func(f string, itemRelationships Relationships, linkRelationshipIDs []RelationshipID) {
		field := func() string {
			return fmt.Sprintf("%s.%s.%s", relatedField, link.ID(), f)
		}
		for _, linkRelationshipID := range linkRelationshipIDs {
			for itemRelationshipID := range itemRelationships {
				if itemRelationshipID == linkRelationshipID {
					alreadyHasRelatedAs = true
				} else {
					updates = append(updates, dal.Update{Field: field(), Value: dal.DeleteField})
				}
			}
		}
	}
	addIfNeeded("relatedAs", relatedItem.RelatedAs, link.RelatedAs)
	addIfNeeded("relatesAs", relatedItem.RelatesAs, link.RelatesAs)

	if alreadyHasRelatedAs {
		if len(v.Related) == 0 {
			v.Related = nil
		}
		return updates, nil
	}

	var relationshipUpdate []dal.Update
	if relationshipUpdate, err = v.AddRelationship(userID, link, now); err != nil {
		return updates, err
	}
	updates = append(updates, relationshipUpdate...)

	return updates, err
}
