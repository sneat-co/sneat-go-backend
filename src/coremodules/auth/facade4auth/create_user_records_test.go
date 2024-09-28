package facade4auth

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/golang/mock/gomock"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/dto4auth"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"testing"
)

func Test_InitUserRecord(t *testing.T) {
	ctx := context.Background()
	type args struct {
		user         facade.UserContext
		userToCreate dto4auth.DataToCreateUser
	}
	tests := []struct {
		name     string
		args     args
		wantUser dbo4userus.UserEntry
		wantErr  bool
	}{
		{
			name: "should_create_user_record",
			args: args{
				user: facade.NewUserContext("test_user_1"),
				userToCreate: dto4auth.DataToCreateUser{
					AuthAccount: appuser.AccountKey{
						Provider: "password",
						ID:       "u1@example.com",
					},
					Names: person.NameFields{
						FirstName: "First",
						LastName:  "UserEntry",
					},
					IanaTimezone: "Europe/Paris",
					Email:        "u1@example.com",
					RemoteClient: dbmodels.RemoteClientInfo{
						HostOrApp:  "unit-test",
						RemoteAddr: "127.0.0.1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// SETUP MOCKS BEGINS

			db := mocks4dal.NewMockDatabase(gomock.NewController(t))
			facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
				return db, nil
			}

			db.EXPECT().RunReadwriteTransaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f dal.RWTxWorker, options ...dal.TransactionOption) error {
				mockCtrl := gomock.NewController(t)
				tx := mocks4dal.NewMockReadwriteTransaction(mockCtrl)
				tx.EXPECT().Get(ctx, gomock.Any()).Return(dal.ErrRecordNotFound).AnyTimes() // TODO: Assert gets
				tx.EXPECT().Insert(ctx, gomock.Any()).Return(nil).AnyTimes()                // TODO: Assert inserts
				return f(ctx, tx)
			})

			sneatauth.GetUserInfo = func(ctx context.Context, uid string) (authUser *sneatauth.AuthUserInfo, err error) {
				authUser = &sneatauth.AuthUserInfo{
					AuthProviderUserInfo: &sneatauth.AuthProviderUserInfo{
						ProviderID: "firebase",
					},
				}
				return
			}
			// SETUP MOCKS ENDS

			// TEST CALL BEGINS
			gotParams, err := CreateUserRecords(ctx, tt.args.user, tt.args.userToCreate)
			// TEST CALL ENDS

			if (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateUserRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.args.userToCreate.Email, gotParams.User.Data.Email)
		})
	}
}
