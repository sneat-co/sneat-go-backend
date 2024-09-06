package facade4spaceus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/golang/mock/gomock"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/person"
	"testing"
)

func TestCreateSpace(t *testing.T) { // TODO: Implement unit tests
	ctx := context.Background()
	user := facade.NewUserContext("TestUser")
	//userKey := dbo4userus.NewUserKey(user.GetID())

	setupMockDb := func() {
		mockCtrl := gomock.NewController(t)
		db := mocks4dal.NewMockDatabase(mockCtrl)
		facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
			return db, nil
		}

		tx := mocks4dal.NewMockReadwriteTransaction(mockCtrl)
		tx.EXPECT().Get(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, record dal.Record) error {
			switch record.Key().Collection() {
			case dbo4userus.UsersCollection:
				record.SetError(nil)
				userDto := record.Data().(*dbo4userus.UserDbo)
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

		facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
			return db, nil
		}
	}

	t.Run("error on bad request", func(t *testing.T) {
		setupMockDb()
		result, err := CreateSpace(ctx, user, dto4spaceus.CreateSpaceRequest{})
		assert.Error(t, err)
		assert.Equal(t, "", result.Space.ID)
		assert.Equal(t, "", result.ContactusSpace.ID)
	})

	t.Run("user's 1st team", func(t *testing.T) {
		setupMockDb()

		result, err := CreateSpace(ctx, user, dto4spaceus.CreateSpaceRequest{Type: core4spaceus.SpaceTypeFamily})
		assert.Nil(t, err)

		assert.NotEqual(t, "", result.Space.ID)
		assert.Nil(t, result.Space.Data.Validate())
		assert.Equal(t, 1, len(result.Space.Data.UserIDs))
		assert.Equal(t, 1, result.Space.Data.Version)
		assert.Equal(t, "contactus", result.ContactusSpace.ID)
	})

}

func Test_getUniqueSpaceID(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	readSession := mocks4dal.NewMockReadSession(mockCtrl)
	readSession.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dal.ErrRecordNotFound)
	teamID, err := getUniqueSpaceID(ctx, readSession, "TestCompany LTD")
	assert.NoError(t, err)
	assert.Equal(t, "testcompanyltd", teamID)
}
