package models4linkage

import (
	"errors"
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

type RelatedItem struct {
	// Brief any // TODO: do we need a brief here?

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

func NewRelatedItem() *RelatedItem {
	return &RelatedItem{
		RelatedAs: make(Relationships, 1),
		RelatesAs: make(Relationships, 1),
	}
}

func (v *RelatedItem) Validate() error {
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
			return errors.New("key is empty string")
		}
		if err := relationshipDetails.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type RelatedByItemID = map[string]*RelatedItem
type RelatedByCollectionID = map[string]RelatedByItemID
type RelatedByModuleID = map[string]RelatedByCollectionID
type RelatedByTeamID = map[string]RelatedByModuleID

const relatedField = "related"

var _ Relatable = (*WithRelatedAndIDs)(nil)

func (v *WithRelatedAndIDs) GetRelated() *WithRelatedAndIDs {
	return v
}

type WithRelated struct {
	// Related defines relationship of the current contact to other contacts. Key is contact ID.
	Related RelatedByTeamID `json:"related,omitempty" firestore:"related,omitempty"`
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
	for teamID, relatedByModuleID := range v.Related {
		if teamID == "" {
			return validation.NewErrBadRecordFieldValue(relatedField, "has empty team ID")
		}
		for moduleID, relatedByCollectionID := range relatedByModuleID {
			if moduleID == "" {
				return validation.NewErrBadRecordFieldValue(
					relatedField+"."+teamID,
					"has empty module ID")
			}
			for collectionID, relatedByItemID := range relatedByCollectionID {
				if collectionID == "" {
					return validation.NewErrBadRecordFieldValue(
						fmt.Sprintf("%s.%s.%s", relatedField, teamID, moduleID),
						"has empty collection ID",
					)
				}
				for itemID, relatedItem := range relatedByItemID {
					if itemID == "" {
						return validation.NewErrBadRecordFieldValue(
							fmt.Sprintf("%s.%s.%s.%s", relatedField, teamID, moduleID, collectionID),
							"has empty item ID")
					}

					key := fmt.Sprintf("%s.%s.%s.%s", teamID, moduleID, collectionID, itemID)
					field := relatedField + "." + key

					if relatedItem == nil {
						return validation.NewErrRecordIsMissingRequiredField(field)
					}

					if err := relatedItem.Validate(); err != nil {
						return validation.NewErrBadRecordFieldValue(field, err.Error())
					}

					if validateID != nil {
						if err := validateID(key); err != nil {
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

	var alreadyHasRelatedAs, alreadyHasRelatesAs bool

	if relatedByModuleID := v.Related[link.TeamID]; relatedByModuleID != nil {
		if relatedByCollectionID := relatedByModuleID[link.ModuleID]; relatedByCollectionID != nil {
			if relatedByItemID := relatedByCollectionID[const4contactus.ContactsCollection]; relatedByItemID != nil {
				if related := relatedByItemID[link.ItemID]; related != nil {
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
					addIfNeeded("relatedAs", related.RelatedAs, link.RelatedAs)
					addIfNeeded("relatesAs", related.RelatesAs, link.RelatesAs)
				}
			}
		}
	}

	if alreadyHasRelatedAs && alreadyHasRelatesAs {
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
