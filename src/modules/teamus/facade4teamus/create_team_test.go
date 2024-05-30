package facade4teamus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/golang/mock/gomock"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/person"
	"testing"
)

func TestCreateTeam(t *testing.T) { // TODO: Implement unit tests
	ctx := context.Background()
	user := facade.NewUser("TestUser")
	//userKey := models4userus.NewUserKey(user.GetID())

	t.Run("error on bad request", func(t *testing.T) {
		response, err := CreateTeam(ctx, user, dto4teamus.CreateTeamRequest{})
		assert.Error(t, err)
		assert.Equal(t, "", response.Team.ID)
	})

	t.Run("user's 1st team", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		db := mocks4dal.NewMockDatabase(mockCtrl)

		tx := mocks4dal.NewMockReadwriteTransaction(mockCtrl)
		tx.EXPECT().Get(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, record dal.Record) error {
			switch record.Key().Collection() {
			case models4userus.UsersCollection:
				record.SetError(nil)
				userDto := record.Data().(*models4userus.UserDbo)
				userDto.CountryID = "--"
				userDto.Status = "active"
				userDto.Gender = dbmodels.GenderMale
				userDto.AgeGroup = dbmodels.AgeGroupAdult
				userDto.Type = briefs4contactus.ContactTypePerson
				userDto.Names = &person.NameFields{
					FirstName: "1st",
					LastName:  "Lastnameoff",
				}
				userDto.Created = dbmodels.CreatedInfo{
					Client: dbmodels.RemoteClientInfo{
						HostOrApp:  "sneat.app",
						RemoteAddr: "127.0.0.1",
					},
				}
				return nil
			default:
				err := dal.ErrRecordNotFound
				record.SetError(err)
				return err
			}
		}).AnyTimes()
		tx.EXPECT().Insert(ctx, gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, record dal.Record, opts ...dal.InsertOption) error {
			return nil
		}).AnyTimes()
		tx.EXPECT().Update(ctx, gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, key *dal.Key, updates []dal.Update, preconditions ...dal.Precondition) error {
			return nil
		}).AnyTimes()
		db.EXPECT().RunReadwriteTransaction(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, worker func(ctx context.Context, tx dal.ReadwriteTransaction) error, o ...dal.TransactionOption) error {
			return worker(ctx, tx)
		}).AnyTimes()

		facade.GetDatabase = func(ctx context.Context) dal.DB {
			return db
		}
		response, err := CreateTeam(ctx, user, dto4teamus.CreateTeamRequest{Type: core4teamus.TeamTypeFamily})
		assert.Nil(t, err)

		assert.NotEqual(t, "", response.Team.ID)
		assert.Nil(t, response.Team.Data.Validate())
		assert.Equal(t, 1, len(response.Team.Data.UserIDs))
		assert.Equal(t, 1, response.Team.Data.Version)
		//assert.Equal(t, 2, len(response.Team.Dbo.UserIDs))

		assert.Nil(t, response.User.Dbo.Validate())
		assert.Equal(t, 1, len(response.User.Dbo.TeamIDs))
		assert.Equal(t, 1, len(response.User.Dbo.Teams))
	})

}

func Test_getUniqueTeamID(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	readSession := mocks4dal.NewMockReadSession(mockCtrl)
	readSession.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dal.ErrRecordNotFound)
	teamID, err := getUniqueTeamID(ctx, readSession, "TestCompany LTD")
	assert.NoError(t, err)
	assert.Equal(t, "testcompanyltd", teamID)
}
