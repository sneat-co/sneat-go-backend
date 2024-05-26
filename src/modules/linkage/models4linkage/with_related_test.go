package models4linkage

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"reflect"
	"testing"
	"time"
)

func TestWithRelatedAndIDs_SetRelationshipToItem(t *testing.T) {
	type fields struct {
		Related    RelatedByModuleID
		relatedIDs []string
	}
	type args struct {
		userID string
		link   Link
		now    time.Time
	}
	now := time.Now()
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantUpdates []dal.Update
	}{
		{
			name:   "set_related_as_parent_for_empty",
			fields: fields{},
			args: args{
				userID: "u1",
				link: Link{
					TeamModuleItemRef: TeamModuleItemRef{
						TeamID:     "team1",
						ModuleID:   const4contactus.ModuleID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     "c2",
					},
					Add: &RolesCommand{
						RolesOfItem: []RelationshipRoleID{"parent"},
					},
				},
				now: now,
			},
			wantUpdates: []dal.Update{
				{
					Field: "related.contactus.contacts", // team1.c2.relatedAs.child
					Value: []*RelatedItem{
						{
							Keys: []RelatedItemKey{
								{TeamID: "team1", ItemID: "c2"},
							},
							RolesOfItem: RelationshipRoles{
								"parent": &RelationshipRole{},
							},
							RolesToItem: RelationshipRoles{
								"child": &RelationshipRole{},
							},
						},
					},
				},
				//{Field: "related.team1.contactus.contacts.c2.relatesAs.child", Value: &RelationshipRole{WithCreatedField: dbmodels.WithCreatedField{Created: dbmodels.Created{By: "u1", On: now.Format(time.DateTime)}}}},
				{Field: "relatedIDs", Value: []string{
					"*",
					"contactus.*",
					"contactus.contacts.*",
					"contactus.contacts.team1.*",
					"contactus.contacts.team1.c2",
				}},
			},
		},
		{
			name:   "set_related_as_child_for_empty",
			fields: fields{},
			args: args{
				userID: "u1",
				link: Link{
					TeamModuleItemRef: TeamModuleItemRef{
						TeamID:     "team1",
						ModuleID:   const4contactus.ModuleID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     "c2",
					},
					Add: &RolesCommand{
						RolesOfItem: []RelationshipRoleID{"child"},
					},
				},
				now: now,
			},
			wantUpdates: []dal.Update{
				{Field: "related.contactus.contacts", // team1.c2.relatedAs.child
					Value: []*RelatedItem{
						{
							Keys: []RelatedItemKey{
								{TeamID: "team1", ItemID: "c2"},
							},
							RolesOfItem: RelationshipRoles{
								"child": &RelationshipRole{},
							},
							RolesToItem: RelationshipRoles{
								"parent": &RelationshipRole{},
							},
						},
					},
				},
				{Field: "relatedIDs", Value: []string{
					"*",
					"contactus.*",
					"contactus.contacts.*",
					"contactus.contacts.team1.*",
					"contactus.contacts.team1.c2",
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &WithRelatedAndIDs{
				WithRelated: WithRelated{
					Related: tt.fields.Related,
				},
				RelatedIDs: tt.fields.relatedIDs,
			}
			gotUpdates, gotErr := v.AddRelationshipAndID(
				tt.args.link,
			)
			if gotErr != nil {
				t.Fatal(gotErr)
			}
			if len(gotUpdates) != len(tt.wantUpdates) {
				t.Errorf("SetRelationshipToItem()\nactual:\n%+v,\nwant:\n%+v", gotUpdates, tt.wantUpdates)
			}
			for i, gotUpdate := range gotUpdates {
				wantUpdate := tt.wantUpdates[i]
				if gotUpdate.Field != wantUpdate.Field {
					t.Errorf("SetRelationshipToItem()[%d]\nactual.Field:\n\t%+v,\nwant.Field:\n\t%+v", i, gotUpdate.Field, wantUpdate.Field)
				}
				if !reflect.DeepEqual(gotUpdate.FieldPath, wantUpdate.FieldPath) {
					t.Errorf("SetRelationshipToItem()[%d]\nactual.Field:\n\t%+v,\nwant.Field:\n\t%+v", i, gotUpdate.FieldPath, wantUpdate.FieldPath)
				}
				if !reflect.DeepEqual(gotUpdate.Value, wantUpdate.Value) {
					t.Errorf("SetRelationshipToItem()[%d]\nactual.Value:\n\t%+v,\nwant.Value:\n\t%+v", i, gotUpdate.Value, wantUpdate.Value)
					if gotUpdate.Field == "related" {
						gotItems, ok := gotUpdate.Value.([]*RelatedItem)
						if !ok {
							t.Errorf("SetRelationshipToItem()[%d]\nactual type:\n\t%T,\nwant type:\n\t%T", i, gotUpdate.Value, wantUpdate.Value)
							return
						}
						wantItems := wantUpdate.Value.([]*RelatedItem)
						if len(gotItems) != len(wantItems) {
							t.Errorf("SetRelationshipToItem()[%d]\nactual.Value:\n\t%+v,\nwant.Value:\n\t%+v", i, gotItems, wantItems)
							return
						}
						for j, gotItem := range gotItems {
							wantItem := wantItems[j]
							if !reflect.DeepEqual(gotItem.Keys, wantItem.Keys) {
								t.Errorf("SetRelationshipToItem()[%d]\nactual.Value[%d].Keys:\n\t%+v,\nwant.Value[%d].Keys:\n\t%+v", i, j, gotItem.Keys, j, wantItem.Keys)
							}
							if !reflect.DeepEqual(gotItem.RolesOfItem, wantItem.RolesOfItem) {
								t.Errorf("SetRelationshipToItem()[%d]\nactual.Value[%d].RolesOfItem:\n\t%+v,\nwant.Value[%d].RolesOfItem:\n\t%+v", i, j, gotItem.RolesOfItem, j, wantItem.RolesOfItem)
							}
							if !reflect.DeepEqual(gotItem.RolesToItem, wantItem.RolesToItem) {
								t.Errorf("SetRelationshipToItem()[%d]\nactual.Value[%d].RolesToItem:\n\t%+v,\nwant.Value[%d].RolesToItem:\n\t%+v", i, j, gotItem.RolesToItem, j, wantItem.RolesToItem)
							}
						}
					}
				}
			}
		})
	}
}
