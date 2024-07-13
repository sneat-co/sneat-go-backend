package dbo4linkage

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

const NoRelatedID = "-"
const AnyRelatedID = "*"

// WithRelatedAndIDs defines relationships of the current contact record to other contacts.
type WithRelatedAndIDs struct {
	WithRelated
	WithRelatedIDs

	//	Example of related field as a JSON and relevant relatedIDs field:
	/*
	   ContactEntry(id="child1") {
	   	relatedIDs: ["space1:parent1:contactus:contacts:parent"],
	   	related: {
	   		"space1": { // Space ID
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

type WithRelatedIDs struct {
	// RelatedIDs holds identifiers of related records - needed for indexed search.
	RelatedIDs []string `json:"relatedIDs,omitempty" firestore:"relatedIDs,omitempty"`
}

func (v *WithRelatedIDs) Validate() error {
	for i, relatedID := range v.RelatedIDs {
		s := strings.TrimSpace(relatedID)
		if s == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "empty contact ID")
		}
		if s != relatedID {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "has leading or trailing spaces")
		}
	}
	return nil
}

func ValidateRelatedAndRelatedIDs(withRelated WithRelated, relatedIDs []string) error {
	if err := withRelated.ValidateRelated(func(relatedID string) error {
		if !slice.Contains(relatedIDs, relatedID) {
			return validation.NewErrBadRecordFieldValue("relatedIDs",
				fmt.Sprintf(`does not have relevant value in 'relatedIDs' field: relatedID="%s"`, relatedID))
		}
		return nil
	}); err != nil {
		return err
	}
	if len(relatedIDs) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("relatedIDs")
	}
	if relatedIDs[0] != AnyRelatedID && relatedIDs[0] != NoRelatedID {
		return validation.NewErrBadRecordFieldValue("relatedIDs[0]", fmt.Sprintf("should be either '%s' or '%s'", AnyRelatedID, NoRelatedID))
	}
	for i, relatedID := range relatedIDs[1:] { // The first item is always either "*" or "-"
		if strings.TrimSpace(relatedID) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "empty contact ID")
		}
		if strings.HasSuffix(relatedID, "."+AnyRelatedID) {
			// TODO: Validate search index values
			continue
		}
		delimitersCount := strings.Count(relatedID, "&")
		switch delimitersCount {
		case 0:
			if relatedID != "*" {
				parts := strings.Split(relatedID, "=")
				if len(parts) != 2 {
					return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "should be in format '{key}={value}', got: "+relatedID)
				}
				switch parts[0] {
				case "m", "s": // Module
					if err := with.ValidateRecordID(parts[1]); err != nil {
						return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), err.Error())
					}
				default:
					return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "single key should start with 'm=' or 's=', got: "+relatedID)
				}
			}
		case 3: // "{module}.{collection}.{space}.{item}"
			relatedRef, err := NewSpaceModuleItemRefFromString(relatedID)
			if err != nil {
				return err
			}

			relatedByCollectionID := withRelated.Related[relatedRef.Module]
			if relatedByCollectionID == nil {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s]' does not have value for module ID=%s", relatedRef.Space, relatedRef.Module))
			}
			relatedItems := relatedByCollectionID[relatedRef.Collection]
			if relatedItems == nil {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s][%s]' does not have value for collection ID=%s", relatedRef.Space, relatedRef.Module, relatedRef.Collection))
			}

			if !HasRelatedItem(relatedItems, RelatedItemKey{SpaceID: relatedRef.Space, ItemID: relatedRef.ItemID}) {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s][%s][%s]' does not have value for item ID=%s", relatedRef.Space, relatedRef.Module, relatedRef.Collection, relatedRef.ItemID))
			}
		}
	}
	return nil
}

// Validate returns error if not valid
func (v *WithRelatedAndIDs) Validate() error {
	if err := v.WithRelatedIDs.Validate(); err != nil {
		return err
	}
	return ValidateRelatedAndRelatedIDs(v.WithRelated, v.RelatedIDs)
}

func (v *WithRelatedAndIDs) AddRelationshipsAndIDs(
	itemRef SpaceModuleItemRef,
	rolesOfItem RelationshipRoles,
	rolesToItem RelationshipRoles, // TODO: needs implementation
) (updates []dal.Update, err error) {
	link := RelationshipItemRolesCommand{}
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
	return v.AddRelationshipAndID(itemRef, link)
	//return nil, errors.New("not implemented yet - AddRelationshipsAndIDs")
}

func UpdateRelatedIDs(withRelated *WithRelated, withRelatedIDs *WithRelatedIDs) (updates []dal.Update) {
	searchIndex := []string{AnyRelatedID}
	withRelatedIDs.RelatedIDs = make([]string, 0)
	for moduleID, relatedByCollectionID := range withRelated.Related {
		searchIndex = append(searchIndex, "m="+moduleID)
		for collectionID, relatedItems := range relatedByCollectionID {
			searchIndex = append(searchIndex, fmt.Sprintf("m=%s&c=%s", moduleID, collectionID))
			spaceIDs := make([]string, 0, len(relatedItems))
			for _, relatedItem := range relatedItems {
				for _, k := range relatedItem.Keys {
					if !slice.Contains(spaceIDs, k.SpaceID) {
						spaceIDs = append(spaceIDs, k.SpaceID)
						searchIndex = append(searchIndex, fmt.Sprintf("s=%s&m=%s&c=%s", k.SpaceID, moduleID, collectionID))
						if spaceKey := fmt.Sprintf("s=%s", k.SpaceID); !slice.Contains(searchIndex, spaceKey) {
							searchIndex = append(searchIndex, spaceKey)
						}
					}
					id := NewSpaceModuleItemRef(k.SpaceID, moduleID, collectionID, k.ItemID).ID()
					withRelatedIDs.RelatedIDs = append(withRelatedIDs.RelatedIDs, id)
				}
			}
		}
	}
	if len(withRelatedIDs.RelatedIDs) == 0 {
		withRelatedIDs.RelatedIDs = []string{NoRelatedID}
		updates = append(updates, dal.Update{Field: "relatedIDs", Value: dal.DeleteField})
	} else {
		withRelatedIDs.RelatedIDs = append(searchIndex, withRelatedIDs.RelatedIDs...)
		updates = append(updates, dal.Update{Field: "relatedIDs", Value: withRelatedIDs.RelatedIDs})
	}
	return
}

func (v *WithRelatedAndIDs) AddRelationshipAndID(
	itemRef SpaceModuleItemRef,
	link RelationshipItemRolesCommand,
) (updates []dal.Update, err error) {
	updates, err = v.WithRelated.AddRelationship(itemRef, link)
	updates = append(updates, UpdateRelatedIDs(&v.WithRelated, &v.WithRelatedIDs)...)
	return
}

func AddRelationshipAndID(
	withRelated *WithRelated,
	withRelatedIDs *WithRelatedIDs,
	itemRef SpaceModuleItemRef,
	link RelationshipItemRolesCommand,
) (updates []dal.Update, err error) {
	updates, err = withRelated.AddRelationship(itemRef, link)
	updates = append(updates, UpdateRelatedIDs(withRelated, withRelatedIDs)...)
	return
}

func RemoveRelatedAndID(withRelated *WithRelated, withRelatedIDs *WithRelatedIDs, ref SpaceModuleItemRef) (updates []dal.Update) {
	updates = withRelated.RemoveRelatedItem(ref)
	updates = append(updates, UpdateRelatedIDs(withRelated, withRelatedIDs)...)
	return updates
}

const (
	RelationshipRoleSpouse   = "spouse"
	RelationshipRoleParent   = "parent"
	RelationshipRoleChild    = "child"
	RelationshipRoleCousin   = "cousin"
	RelationshipRoleSibling  = "sibling"
	RelationshipRolePartner  = "partner"
	RelationshipRoleTeammate = "space-mate"
)

// Should provide a way for modules to register opposite roles?
var oppositeRoles = map[RelationshipRoleID]RelationshipRoleID{
	RelationshipRoleParent: RelationshipRoleChild,
	RelationshipRoleChild:  RelationshipRoleParent,
}

// Should provide a way for modules to register reciprocal roles?
var reciprocalRoles = []string{
	RelationshipRoleSpouse,
	RelationshipRoleSibling,
	RelationshipRoleCousin,
	RelationshipRolePartner,
	RelationshipRoleTeammate,
}

func IsReciprocalRole(role RelationshipRoleID) bool {
	return slice.Contains(reciprocalRoles, role)
}

// GetOppositeRole returns relationship ID for the opposite direction
func GetOppositeRole(relationshipRoleID RelationshipRoleID) RelationshipRoleID {
	if relationshipRoleID == "" {
		return ""
	}
	if IsReciprocalRole(relationshipRoleID) {
		return relationshipRoleID
	}
	return oppositeRoles[relationshipRoleID]
}
