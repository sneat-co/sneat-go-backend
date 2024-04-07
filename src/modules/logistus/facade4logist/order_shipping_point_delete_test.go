package facade4logist

import (
	"context"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/mocks4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_txDeleteShippingPoint(t *testing.T) {
	ctx := context.Background()

	type args struct {
		params  *OrderWorkerParams
		request dto4logist.OrderShippingPointRequest
	}

	sharedPreAssert := func(t *testing.T, args args) {
		assert.False(t, args.params.Changed.HasChanges())
		assert.Equal(t, mocks4logist.Dto1ShippingPointsCount, len(args.params.Order.Dto.ShippingPoints), fmt.Sprintf("%+v", args.params.Order.Dto.ShippingPoints))
		assert.Equal(t, 1, len(args.params.Order.Dto.Segments), fmt.Sprintf("%+v", args.params.Order.Dto.Segments))
	}
	sharedPostAssert := func(t *testing.T, err error, args args) {
		assert.True(t, args.params.Changed.HasChanges())
		assert.Equal(t, mocks4logist.Dto1ShippingPointsCount-1, len(args.params.Order.Dto.ShippingPoints), fmt.Sprintf("%+v", args.params.Order.Dto.ShippingPoints))
		assert.Equal(t, 0, len(args.params.Order.Dto.Segments), fmt.Sprintf("%+v", args.params.Order.Dto.Segments))
	}
	tests := []struct {
		name       string
		args       args
		preAssert  func(t *testing.T, args args)
		postAssert func(t *testing.T, err error, args args)
	}{
		{
			name: "should_pass",
			args: args{
				params: &OrderWorkerParams{
					Order: models4logist.NewOrderWithData("team1", "order1", mocks4logist.ValidOrderDto1(t)),
				},
				request: dto4logist.OrderShippingPointRequest{
					ShippingPointID: mocks4logist.ShippingPoint1WithSingleContainerID,
					OrderRequest:    dto4logist.NewOrderRequest("team1", "order1"),
				},
			},
			preAssert: func(t *testing.T, args args) {
				t.Helper()
				sharedPreAssert(t, args)
			},
			postAssert: func(t *testing.T, err error, args args) {
				t.Helper()
				assert.Nil(t, err)
				sharedPostAssert(t, err, args)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := mocks4logist.MockTx(t)
			tt.preAssert(t, tt.args)
			err := txDeleteShippingPoint(ctx, tx, tt.args.params, tt.args.request)
			tt.postAssert(t, err, tt.args)
		})
	}
}
