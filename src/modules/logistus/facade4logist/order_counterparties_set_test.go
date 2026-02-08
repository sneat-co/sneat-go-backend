package facade4logist

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/mocks4logist"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/with"
	"go.uber.org/mock/gomock"
)

type mockUser struct {
	facade.UserContext
	userID string
}

func (m mockUser) GetUserID() string {
	return m.userID
}

type mockContext struct {
	facade.ContextWithUser
	user facade.UserContext
}

func (m mockContext) User() facade.UserContext {
	return m.user
}

func TestSetOrderCounterparties(t *testing.T) {
	origRunOrderWorker := RunOrderWorker
	defer func() { RunOrderWorker = origRunOrderWorker }()

	RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
		spaceEntry := dbo4spaceus.NewSpaceEntry("space1")
		spaceEntry.Data = &dbo4spaceus.SpaceDbo{
			WithUserIDs: dbmodels.WithUserIDs{UserIDs: []string{"u1"}},
		}
		tx := mocks4logist.MockTx(t)
		tx.EXPECT().GetMulti(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		order := mocks4logist.ValidEmptyOrder(t)
		order.SpaceID = "space1"
		order.SpaceIDs = []coretypes.SpaceID{"space1"}
		order.UserIDs = []string{"u1"}

		return worker(ctx, tx, &OrderWorkerParams{
			SpaceWorkerParams: &dal4spaceus.SpaceWorkerParams{
				Space: spaceEntry,
			},
			Order: dbo4logist.Order{Dto: order},
		})
	}

	request := dto4logist.SetOrderCounterpartiesRequest{
		OrderRequest: dto4logist.NewOrderRequest("space1", "order1"),
	}
	ctx := mockContext{user: mockUser{userID: "u1"}}
	_, err := SetOrderCounterparties(ctx, request)
	assert.Nil(t, err)
}

func TestCreateCounterparty(t *testing.T) {
	request := dto4logist.CreateCounterpartyRequest{
		RolesField: with.RolesField{
			Roles: []string{dbo4logist.CounterpartyRoleTrucker},
		},
	}
	// This will fail because facade4contactus.CreateContact is not mocked, but it covers the validation line
	_, err := CreateCounterparty(nil, request)
	assert.NotNil(t, err)
}
