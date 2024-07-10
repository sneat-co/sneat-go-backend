package facade4logist

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/mocks4logist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteContainerPoints(t *testing.T) {
	ctx := context.Background()

	type args struct {
		request dto4logist.ContainerPointsRequest
		params  OrderWorkerParams
	}

	sharedPreAssert := func(t *testing.T, args args) {
		order := args.params.Order
		assert.Equal(t, mocks4logist.Dto1ContainerPointsCount, len(order.Dto.ContainerPoints))
		assert.Equal(t, 1, len(order.Dto.Segments))
	}

	sharedPostAssert := func(t *testing.T, args args, err error) {
		for _, shippingPointID := range args.request.ShippingPointIDs {
			assert.Nil(t, args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, shippingPointID))
		}
	}

	tests := []struct {
		name       string
		args       args
		preAssert  func(t *testing.T, args args)
		postAssert func(t *testing.T, args args, err error)
	}{
		{
			name: "delete_from_shipping_point_with_single_container",
			args: args{
				request: dto4logist.ContainerPointsRequest{
					ContainerID:      mocks4logist.Container2ID,
					ShippingPointIDs: []string{mocks4logist.ShippingPoint1WithSingleContainerID},
				},
				params: OrderWorkerParams{
					Order: dbo4logist.NewOrderWithData("space1", "order1",
						mocks4logist.ValidOrderDto1(t)),
				},
			},
			preAssert: func(t *testing.T, args args) {
				sharedPreAssert(t, args)
				for _, shippingPointID := range args.request.ShippingPointIDs {
					assert.NotNil(t, args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, shippingPointID))
				}
			},
			postAssert: func(t *testing.T, args args, err error) {
				sharedPostAssert(t, args, err)
				assert.Equal(t, 4, len(args.params.Order.Dto.ContainerPoints))
			},
		},
		{
			name: "delete_from_shipping_point_with_few_containers",
			args: args{
				request: dto4logist.ContainerPointsRequest{
					ContainerID:      mocks4logist.Container1ID,
					ShippingPointIDs: []string{mocks4logist.ShippingPoint2With2ContainersID},
				},
				params: OrderWorkerParams{
					Order: dbo4logist.NewOrderWithData("space1", "order1",
						mocks4logist.ValidOrderDto1(t)),
				},
			},
			preAssert: func(t *testing.T, args args) {
				sharedPreAssert(t, args)
				for _, shippingPointID := range args.request.ShippingPointIDs {
					assert.NotNil(t, args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, shippingPointID))
				}
			},
			postAssert: func(t *testing.T, args args, err error) {
				sharedPostAssert(t, args, err)
				assert.Equal(t, 4, len(args.params.Order.Dto.ContainerPoints))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preAssert(t, tt.args)
			err := txDeleteContainerPoints(ctx, nil, &tt.args.params, tt.args.request)
			tt.postAssert(t, tt.args, err)
		})
	}
}
