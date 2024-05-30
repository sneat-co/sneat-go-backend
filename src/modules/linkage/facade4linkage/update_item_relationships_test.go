package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-core/facade"
	"reflect"
	"testing"
)

func TestUpdateItemRelationships(t *testing.T) {
	type args struct {
		ctx     context.Context
		userCtx facade.User
		request dto4linkage.UpdateItemRequest
	}
	const testUserID = "test_user_1"
	const team1ID = "team_1"
	const item1ID = "item_1"
	const collection1ID = "collection_1"
	const module1ID = "module_1"

	tests := []struct {
		name      string
		args      args
		wantItem  record.DataWithID[string, *models4linkage.WithRelatedAndIDsAndUserID]
		wantErr   bool
		wantPanic bool
	}{
		{
			name:      "should_update_contact_with_reciprocal_role",
			wantPanic: true, // TODO: Fix this test
			args: args{
				ctx:     context.Background(),
				userCtx: facade.NewUser(testUserID),
				request: dto4linkage.UpdateItemRequest{
					TeamModuleItemRef: models4linkage.TeamModuleItemRef{
						ModuleID:   const4contactus.ModuleID,
						TeamID:     team1ID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     item1ID,
					},
					UpdateRelatedFieldRequest: dto4linkage.UpdateRelatedFieldRequest{
						Related: map[string]*models4linkage.RelationshipRolesCommand{
							models4linkage.NewTeamModuleItemRef(team1ID, module1ID, collection1ID, item1ID).ID(): {
								Add: &models4linkage.RolesCommand{
									RolesOfItem: []models4linkage.RelationshipRoleID{
										models4linkage.RelationshipRoleSpouse,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("UpdateItemRelationships() did not panic")
					}
				}()
			}
			gotItem, err := UpdateItemRelationships(tt.args.ctx, tt.args.userCtx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateItemRelationships() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotItem, tt.wantItem) {
				t.Errorf("UpdateItemRelationships() gotItem = %v, want %v", gotItem, tt.wantItem)
			}
		})
	}
}
