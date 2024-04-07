package facade4logist

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/mocks4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_addOrderShippingPointTx(t *testing.T) {
	type args struct {
		request dto4logist.AddOrderShippingPointRequest
		params  *OrderWorkerParams
	}

	ctx := context.Background()
	tests := []struct {
		name              string
		args              args
		wantShippingPoint *models4logist.OrderShippingPoint
		assertPost        func(t *testing.T, err error, args args)
	}{
		{
			name: "OK",
			args: args{
				params: &OrderWorkerParams{
					Order: models4logist.NewOrderWithData("team1", "order1", mocks4logist.ValidOrderDto1(t)),
				},
				request: dto4logist.AddOrderShippingPointRequest{
					LocationContactID: mocks4logist.Dispatcher1warehouse1ContactID,
					Tasks:             []models4logist.ShippingPointTask{"load"},
					Containers: []dto4logist.AddContainerPoint{
						{ID: mocks4logist.Container1ID, Tasks: []models4logist.ShippingPointTask{models4logist.ShippingPointTaskLoad}},
						{ID: mocks4logist.Container2ID, Tasks: []models4logist.ShippingPointTask{models4logist.ShippingPointTaskUnload}},
					},
					OrderRequest: dto4logist.NewOrderRequest("order1", "team1"),
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
