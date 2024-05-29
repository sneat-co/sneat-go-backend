package models4linkage

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
)

const NoRelatedID = "-"
const AnyRelatedID = "*"

// WithRelatedAndIDs defines relationships of the current contact record to other contacts.
type WithRelatedAndIDs struct {
	WithRelated

	// RelatedIDs holds identifiers of related records - needed for indexed search.
	RelatedIDs []string `json:"relatedIDs,omitempty" firestore:"relatedIDs,omitempty"`

	//	Example of related field as a JSON and relevant relatedIDs field:
	/*
	   Contact(id="child1") {
	   	relatedIDs: ["team1:parent1:contactus:contacts:parent"],
	   	related: {
	   		"team1": { // Team ID
	   			"contactus": { // Module ID
	   				"contacts": { // Collection
	   					"parent1": { // Item ID
	   						relatedAs: {
	   							"parent": {} // RelationshipRole ID
	   						}
	   						relatesAs: {
	   							"child": {} // RelationshipRole ID
	   						},
	   					},
	   				}
	   			},
	   		},
	   	}
	   }
	*/
}

// Validate returns error if not valid
func (v *WithRelatedAndIDs) Validate() error {
	if err := v.ValidateRelated(func(relatedID string) error {
		if !slice.Contains(v.RelatedIDs, relatedID) {
			return validation.NewErrBadRecordFieldValue("relatedIDs",
				fmt.Sprintf(`does not have relevant value in 'relatedIDs' field: relatedID="%s"`, relatedID))
		}
		return nil
	}); err != nil {
		return err
	}
	if len(v.RelatedIDs) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("relatedIDs")
	}
	if v.RelatedIDs[0] != AnyRelatedID && v.RelatedIDs[0] != NoRelatedID {
		return validation.NewErrBadRecordFieldValue("relatedIDs[0]", fmt.Sprintf("should be either '%s' or '%s'", AnyRelatedID, NoRelatedID))
	}
	for i, relatedID := range v.RelatedIDs[1:] { // The first item is always either "*" or "-"
		if strings.TrimSpace(relatedID) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "empty contact ID")
		}
		if strings.HasSuffix(relatedID, "."+AnyRelatedID) {
			// TODO: Validate search index values
			continue
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
	relatedTo TeamModuleItemRef,
	rolesOfItem RelationshipRoles,
	rolesToItem RelationshipRoles, // TODO: needs implementation
) (updates []dal.Update, err error) {
	link := Link{
		TeamModuleItemRef: relatedTo,
	}
	if len(rolesOfItem) > 0 {
		if link.Add == nil {
			link.Add = new(RolesCommand)
		}
		for roleOfItem := range rolesOfItem {
			link.Add.RolesOfItem = append(link.Add.RolesOfItem, roleOfItem)
		}
	}
	if len(rolesToItem) > 0 {
		if link.Remove == nil {
			link.Remove = new(RolesCommand)
		}
		for roleToItem := range rolesToItem {
			link.Remove.RolesToItem = append(link.Remove.RolesToItem, roleToItem)
		}
	}
	return v.AddRelationshipAndID(link)
	//return nil, errors.New("not implemented yet - AddRelationshipsAndIDs")
}

func (v *WithRelatedAndIDs) UpdateRelatedIDs() (updates []dal.Update) {
	searchIndex := []string{AnyRelatedID}
	v.RelatedIDs = make([]string, 0)
	for moduleID, relatedByCollectionID := range v.Related {
		searchIndex = append(searchIndex, fmt.Sprintf("%s.%s", moduleID, AnyRelatedID))
		for collectionID, relatedItems := range relatedByCollectionID {
			searchIndex = append(searchIndex, fmt.Sprintf("%s.%s.%s", moduleID, collectionID, AnyRelatedID))
			teamIDs := make([]string, 0, len(relatedItems))
			for _, relatedItem := range relatedItems {
				for _, k := range relatedItem.Keys {
					if !slice.Contains(teamIDs, k.TeamID) {
						teamIDs = append(teamIDs, k.TeamID)
						searchIndex = append(searchIndex, fmt.Sprintf("%s.%s.%s.%s", moduleID, collectionID, k.TeamID, AnyRelatedID))
					}
					id := NewTeamModuleDocRef(k.TeamID, moduleID, collectionID, k.ItemID).ID()
					v.RelatedIDs = append(v.RelatedIDs, id)
				}
			}
		}
	}
	if len(v.RelatedIDs) == 0 {
		v.RelatedIDs = []string{NoRelatedID}
		updates = append(updates, dal.Update{Field: "relatedIDs", Value: dal.DeleteField})
	} else {
		v.RelatedIDs = append(searchIndex, v.RelatedIDs...)
		updates = append(updates, dal.Update{Field: "relatedIDs", Value: v.RelatedIDs})
	}
	return
}

func (v *WithRelatedAndIDs) AddRelationshipAndID(
	link Link,
) (updates []dal.Update, err error) {
	updates, err = v.WithRelated.AddRelationship(link)
	updates = append(updates, v.UpdateRelatedIDs()...)
	return
}

func (v *WithRelatedAndIDs) RemoveRelatedAndID(ref TeamModuleItemRef) (updates []dal.Update) {
	updates = v.WithRelated.RemoveRelatedItem(ref)
	updates = append(updates, v.UpdateRelatedIDs()...)
	return updates
}

// GetOppositeRole returns relationship ID for the opposite direction
func GetOppositeRole(relationshipRoleID RelationshipRoleID) RelationshipRoleID {
	// TODO: Move to contactus module as this relationships are relevant to contacts only?
	switch relationshipRoleID {
	case "sibling", "spouse", "partner", "team-mate":
		return relationshipRoleID
	case "parent":
		return "child"
	case "child":
		return "parent"
	}
	return ""
}
