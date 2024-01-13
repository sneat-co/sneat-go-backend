package models4linkage

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// WithRelatedAndIDs defines relationship of the current contact record to other contacts.
type WithRelatedAndIDs struct {
	/* Example of related field as a JSON:

	Contact(id="child1") {
		relatedIDs: ["team1:parent1:contactus:contacts:parent"],
		related: {
			"team1": { // Team ID
				"contactus": { // Module ID
					"contacts": { // Collection
						"parent1": { // Item ID
							relatedAs: {
								"parent": {} // Relationship ID
							}
							relatesAs: {
								"child": {} // Relationship ID
							},
						},
					}
				},
			},
		}
	}
	*/

	WithRelated

	// RelatedIDs is a list of IDs of records that are related to the current record - this is needed for indexed search.
	RelatedIDs []string `json:"relatedIDs,omitempty" firestore:"relatedIDs,omitempty"`
}

// Validate returns error if not valid
func (v *WithRelatedAndIDs) Validate() error {
	if err := v.ValidateRelated(func(relatedID string) error {
		if !slice.Contains(v.RelatedIDs, relatedID) {
			return validation.NewErrBadRecordFieldValue("relatedIDs",
				"does not have relevant value in 'relatedIDs' field: "+relatedID)
		}
		return nil
	}); err != nil {
		return err
	}
	for i, relatedID := range v.RelatedIDs {
		if strings.TrimSpace(relatedID) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "empty contact ID")
		}
		relatedRef := NewTeamModuleDocRefFromString(relatedID)

		relatedByModuleID := v.Related[relatedRef.TeamID]
		if relatedByModuleID == nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related'  does not have value for team ID=%s", relatedRef.TeamID))
		}
		relatedByCollectionID := relatedByModuleID[relatedRef.ModuleID]
		if relatedByCollectionID == nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s]' does not have value for module ID=%s", relatedRef.TeamID, relatedRef.ModuleID))
		}
		relatedByItemID := relatedByCollectionID[relatedRef.Collection]
		if relatedByItemID == nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s][%s]' does not have value for collection ID=%s", relatedRef.TeamID, relatedRef.ModuleID, relatedRef.Collection))
		}
		_, ok := relatedByItemID[relatedRef.ItemID]
		if !ok {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s][%s][%s]' does not have value for item ID=%s", relatedRef.TeamID, relatedRef.ModuleID, relatedRef.Collection, relatedRef.ItemID))
		}
	}
	return nil
}

func (v *WithRelatedAndIDs) SetRelationshipsToItem(
	userID string,
	recordRef TeamModuleDocRef,
	relatedTo TeamModuleDocRef,
	relatedAs Relationships,
	relatesAs Relationships,
	now time.Time,
) (updates []dal.Update, err error) {
	link := Link{
		TeamModuleDocRef: relatedTo,
	}
	for relatedAsID := range relatedAs {
		link.RelatesAs = append(link.RelatesAs, relatedAsID)
	}
	return v.SetRelationshipToItem(userID, link, now)
	//return nil, errors.New("not implemented yet - SetRelationshipsToItem")
}

func (v *WithRelatedAndIDs) SetRelationshipToItem(
	userID string,
	link Link,
	now time.Time,
) (updates []dal.Update, err error) {
	updates, err = v.WithRelated.SetRelationshipToItem(userID, link, now)
	v.updateRelatedIDs()
	updates = append(updates, dal.Update{Field: "relatedIDs", Value: v.RelatedIDs})
	return
}

func (v *WithRelatedAndIDs) updateRelatedIDs() {
	v.RelatedIDs = make([]string, 0, len(v.Related))
	for teamID, relatedByModuleID := range v.Related {
		for moduleID, relatedByCollectionID := range relatedByModuleID {
			for collectionID, relatedByItemID := range relatedByCollectionID {
				for itemID := range relatedByItemID {
					id := NewTeamModuleDocRef(teamID, moduleID, collectionID, itemID).ID()
					v.RelatedIDs = append(v.RelatedIDs, id)
				}
			}
		}
	}
}

func (v *WithRelatedAndIDs) AddRelationship(
	userID string,
	link Link,
	now time.Time,
) (updates []dal.Update, err error) {
	if updates, err = v.WithRelated.AddRelationship(userID, link, now); err != nil {
		return nil, err
	}
	relatedItemID := link.TeamModuleDocRef.ID()
	if !slice.Contains(v.RelatedIDs, relatedItemID) {
		v.RelatedIDs = append(v.RelatedIDs, relatedItemID)
		updates = append(updates, dal.Update{Field: "relatedIDs", Value: v.RelatedIDs})
	}
	return
}

func (v *WithRelated) AddRelationship(
	userID string,
	link Link,
	now time.Time,
) (updates []dal.Update, err error) {
	if err := link.Validate(); err != nil {
		return nil, err
	}
	if v.Related == nil {
		v.Related = make(RelatedByTeamID, 1)
	}

	for _, linkRelatedAs := range link.RelatesAs {
		if relatesAs := GetRelatesAsFromRelated(linkRelatedAs); relatesAs != "" && !slice.Contains(link.RelatesAs, relatesAs) {
			link.RelatesAs = append(link.RelatesAs, "child")
		}
	}

	relatedByModuleID := v.Related[link.TeamID]
	if relatedByModuleID == nil {
		relatedByModuleID = make(RelatedByModuleID, 1)
		v.Related[link.TeamID] = relatedByModuleID
	}

	relatedByCollectionID := relatedByModuleID[link.ModuleID]
	if relatedByCollectionID == nil {
		relatedByCollectionID = make(RelatedByCollectionID, 1)
		relatedByModuleID[const4contactus.ModuleID] = relatedByCollectionID
	}

	relatedByItemID := relatedByCollectionID[const4contactus.ContactsCollection]
	if relatedByItemID == nil {
		relatedByItemID = make(RelatedByItemID, 1)
		relatedByCollectionID[const4contactus.ContactsCollection] = relatedByItemID
	}

	relatedItem := relatedByItemID[link.ItemID]
	if relatedItem == nil {
		relatedItem = NewRelatedItem()
		relatedByItemID[link.ItemID] = relatedItem
	}

	addRelationship := func(field string, relationshipIDs []RelationshipID, relationships Relationships) Relationships {
		if relationships == nil {
			relationships = make(Relationships, len(relationshipIDs))
		}
		for _, relationshipID := range relationshipIDs {
			if relationship := relationships[relationshipID]; relationship == nil {
				relationship = &Relationship{
					CreatedField: with.CreatedField{
						Created: with.Created{
							By: userID,
							At: now.Format(time.DateOnly),
						},
					},
				}
				relationships[relationshipID] = relationship
				updates = append(updates, dal.Update{
					Field: fmt.Sprintf("related.%s.%s.%s", link.ID(), field, relationshipID),
					Value: relationship,
				})
			}
		}
		return relationships
	}

	relatedItem.RelatedAs = addRelationship("relatedAs", link.RelatedAs, relatedItem.RelatedAs)
	relatedItem.RelatedAs = addRelationship("relatesAs", link.RelatesAs, relatedItem.RelatesAs)

	return updates, nil
}

func (v *WithRelatedAndIDs) RemoveRelationship(ref TeamModuleDocRef) (updates []dal.Update) {
	updates = v.WithRelated.RemoveRelationshipToContact(ref)
	if slice.Contains(v.RelatedIDs, ref.ID()) {
		v.RelatedIDs = slice.RemoveInPlace(ref.ID(), v.RelatedIDs)
		updates = append(updates, dal.Update{
			Field: "relatedIDs",
			Value: v.RelatedIDs,
		})
	}
	return updates
}

// GetRelatesAsFromRelated returns relationship ID for the opposite direction
// TODO: Move to contactus module as relationships are not dedicated to contacts only?
func GetRelatesAsFromRelated(relatedAs RelationshipID) RelationshipID {
	switch relatedAs {
	case "parent":
		return "child"
	case "spouse":
		return "spouse"
	}
	return ""
}
