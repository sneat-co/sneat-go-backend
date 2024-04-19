package models4linkage

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/strongo/strongoapp/with"
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
		//recordRef TeamModuleDocRef
		link Link
		now  time.Time
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
				//recordRef: TeamModuleDocRef{
				//	TeamID:     "team1",
				//	ModuleID:   const4contactus.ModuleID,
				//	Collection: const4contactus.ContactsCollection,
				//	ItemID:     "u1c",
				//},
				link: Link{
					TeamModuleDocRef: TeamModuleDocRef{
						TeamID:     "team1",
						ModuleID:   const4contactus.ModuleID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     "c2",
					},
					RelatedAs: []RelationshipID{"parent"},
				},
				now: now,
			},
			wantUpdates: []dal.Update{
				{Field: "related.contactus.contacts.team1.c2.relatedAs.parent", Value: &Relationship{
					CreatedField: with.CreatedField{
						Created: with.Created{By: "u1", At: now.Format(time.DateOnly)}}},
				},
				//{Field: "related.team1.contactus.contacts.c2.relatesAs.child", Value: &Relationship{WithCreatedField: dbmodels.WithCreatedField{Created: dbmodels.Created{By: "u1", On: now.Format(time.DateOnly)}}}},
				{Field: "relatedIDs", Value: []string{"contactus.contacts.team1.c2"}},
			},
		},
		{
			name:   "set_related_as_child_for_empty",
			fields: fields{},
			args: args{
				userID: "u1",
				//recordRef: TeamModuleDocRef{
				//	TeamID:     "team1",
				//	ModuleID:   const4contactus.ModuleID,
				//	Collection: const4contactus.ContactsCollection,
				//	ItemID:     "u1c",
				//},
				link: Link{
					TeamModuleDocRef: TeamModuleDocRef{
						TeamID:     "team1",
						ModuleID:   const4contactus.ModuleID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     "c2",
					},
					RelatedAs: []RelationshipID{"child"},
				},
				now: now,
			},
			wantUpdates: []dal.Update{
				{Field: "related.contactus.contacts.team1.c2.relatedAs.child",
					Value: &Relationship{
						CreatedField: with.CreatedField{
							Created: with.Created{By: "u1",
								At: now.Format(time.DateOnly),
							},
						},
					},
				},
				{Field: "relatedIDs", Value: []string{"contactus.contacts.team1.c2"}},
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
			gotUpdates, gotErr := v.SetRelationshipToItem(
				tt.args.userID,
				tt.args.link,
				tt.args.now,
			)
			if gotErr != nil {
				t.Fatal(gotErr)
			}
			if !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
				t.Errorf("SetRelationshipToItem()\nactual:\n%+v,\nwant:\n%+v", gotUpdates, tt.wantUpdates)
			}
		})
	}
}
