package models4linkage

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
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

		relatedByCollectionID := v.Related[relatedRef.ModuleID]
		if relatedByCollectionID == nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s]' does not have value for module ID=%s", relatedRef.TeamID, relatedRef.ModuleID))
		}
		relatedItems := relatedByCollectionID[relatedRef.Collection]
		if relatedItems == nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s][%s]' does not have value for collection ID=%s", relatedRef.TeamID, relatedRef.ModuleID, relatedRef.Collection))
		}

		if !HasRelatedItem(relatedItems, RelatedItemKey{TeamID: relatedRef.TeamID, ItemID: relatedRef.ItemID}) {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s][%s][%s]' does not have value for item ID=%s", relatedRef.TeamID, relatedRef.ModuleID, relatedRef.Collection, relatedRef.ItemID))
		}
	}
	return nil
}

func (v *WithRelatedAndIDs) AddRelationshipsAndIDs(
	/*recordRef*/ _ TeamModuleItemRef, // TODO: handle or remove
	relatedTo TeamModuleItemRef,
	relatedAs Relationships,
	/*relatesAs*/ _ Relationships, // TODO: needs implementation
) (updates []dal.Update, err error) {
	link := Link{
		TeamModuleItemRef: relatedTo,
	}
	for relatedAsID := range relatedAs {
		link.RelatesAs = append(link.RelatesAs, relatedAsID)
	}
	return v.AddRelationshipAndID(link)
	//return nil, errors.New("not implemented yet - AddRelationshipsAndIDs")
}

func (v *WithRelatedAndIDs) updateRelatedIDs() (updates []dal.Update) {
	v.RelatedIDs = make([]string, 0, len(v.Related))
	for moduleID, relatedByCollectionID := range v.Related {
		for collectionID, relatedItems := range relatedByCollectionID {
			for _, relatedItem := range relatedItems {
				for _, k := range relatedItem.Keys {
					id := NewTeamModuleDocRef(k.TeamID, moduleID, collectionID, k.ItemID).ID()
					v.RelatedIDs = append(v.RelatedIDs, id)
				}
			}
		}
	}
	updates = append(updates, dal.Update{Field: "relatedIDs", Value: v.RelatedIDs})
	return
}

func (v *WithRelatedAndIDs) AddRelationshipAndID(
	link Link,
) (updates []dal.Update, err error) {
	updates, err = v.WithRelated.AddRelationship(link)
	updates = append(updates, v.updateRelatedIDs()...)
	return
}

func (v *WithRelatedAndIDs) RemoveRelatedAndID(ref TeamModuleItemRef) (updates []dal.Update) {
	updates = v.WithRelated.RemoveRelatedItem(ref)
	updates = append(updates, v.updateRelatedIDs()...)
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
