package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/const4contactus"
	dbo4linkage2 "github.com/sneat-co/sneat-go-backend/src/coremodules/linkage/dbo4linkage"
	dto4linkage2 "github.com/sneat-co/sneat-go-backend/src/coremodules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-go-core/facade"
	"reflect"
	"testing"
)

func TestUpdateItemRelationships(t *testing.T) {
	type args struct {
		ctx     context.Context
		userCtx facade.UserContext
		request dto4linkage2.UpdateItemRequest
	}
	const testUserID = "test_user_1"
	const space1ID = "space_1"
	const item1ID = "item_1"
	const collection1ID = "collection_1"
	const module1ID = "module_1"

	facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
		return nil, nil
	}

	tests := []struct {
		name      string
		args      args
		wantItem  record.DataWithID[string, *dbo4linkage2.WithRelatedAndIDsAndUserID]
		wantErr   bool
		wantPanic bool
	}{
		{
			name:      "should_update_contact_with_reciprocal_role",
			wantPanic: true, // TODO: Fix this test
			args: args{
				ctx:     context.Background(),
				userCtx: facade.NewUserContext(testUserID),
				request: dto4linkage2.UpdateItemRequest{
					SpaceModuleItemRef: dbo4linkage2.SpaceModuleItemRef{
						Module:     const4contactus.ModuleID,
						Space:      space1ID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     item1ID,
					},
					UpdateRelatedFieldRequest: dto4linkage2.UpdateRelatedFieldRequest{
						Related: []dbo4linkage2.RelationshipItemRolesCommand{
							{
								ItemRef: dbo4linkage2.NewSpaceModuleItemRef(space1ID, module1ID, collection1ID, item1ID),
								Add: &dbo4linkage2.RolesCommand{
									RolesOfItem: []dbo4linkage2.RelationshipRoleID{
										dbo4linkage2.RelationshipRoleSpouse,
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
