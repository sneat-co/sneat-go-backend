package facade4logist

import (
	"context"
	"testing"

	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/mocks4logist"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

func TestAddOrderShippingPoint(t *testing.T) {
	origRunOrderWorker := RunOrderWorker
	defer func() { RunOrderWorker = origRunOrderWorker }()

	RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
		spaceEntry := dbo4spaceus.NewSpaceEntry("space1")
		spaceEntry.Data = &dbo4spaceus.SpaceDbo{
			WithUserIDs: dbmodels.WithUserIDs{UserIDs: []string{"u1"}},
		}

		tx := mocks4logist.MockTx(t)
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

	request := dto4logist.AddOrderShippingPointRequest{
		OrderRequest:      dto4logist.NewOrderRequest("space1", "order1"),
		LocationContactID: mocks4logist.Dispatcher1warehouse1ContactID,
		Tasks:             []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskLoad},
	}
	_, err := AddOrderShippingPoint(nil, request)
	assert.Nil(t, err)
}

func Test_addOrderShippingPointTx(t *testing.T) {
	type args struct {
		request dto4logist.AddOrderShippingPointRequest
		params  *OrderWorkerParams
	}

	ctx := context.Background()
	tests := []struct {
		name              string
		args              args
		wantShippingPoint *dbo4logist.OrderShippingPoint
		assertPost        func(t *testing.T, err error, args args)
	}{
		{
			name: "OK",
			args: args{
				params: &OrderWorkerParams{
					SpaceWorkerParams: &dal4spaceus.SpaceWorkerParams{
						Space: dbo4spaceus.NewSpaceEntry("space1"),
					},
					Order: dbo4logist.NewOrderWithData("space1", "order1", mocks4logist.ValidOrderDto1(t)),
				},
				request: dto4logist.AddOrderShippingPointRequest{
					LocationContactID: mocks4logist.Dispatcher1warehouse1ContactID,
					Tasks:             []dbo4logist.ShippingPointTask{"load"},
					Containers: []dto4logist.AddContainerPoint{
						{ID: mocks4logist.Container1ID, Tasks: []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskLoad}},
						{ID: mocks4logist.Container2ID, Tasks: []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskUnload}},
					},
					OrderRequest: dto4logist.NewOrderRequest("order1", "space1"),
				},
			},
			assertPost: func(t *testing.T, err error, args args) {
				if err != nil {
					t.Errorf("addOrderShippingPointTx() error = %v, wantErr %v", err, nil)
				}
				assert.Nil(t, args.params.Order.Dto.Validate(), "order must be valid")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := mocks4logist.MockTx(t)
			gotShippingPoint, err := addOrderShippingPointTx(ctx, tx, tt.args.request, tt.args.params)
			tt.assertPost(t, err, tt.args)
			assert.NotNil(t, gotShippingPoint)
		})
	}
}
