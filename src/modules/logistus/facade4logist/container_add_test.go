package facade4logist

import (
	"strings"
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/mocks4logist"
	"github.com/stretchr/testify/assert"
)

func Test_addContainersTx(t *testing.T) {
	type args struct {
		params  *OrderWorkerParams
		request dto4logist.AddContainersRequest
	}
	sharedPreAssert := func(t *testing.T, args args) {
		if !assert.Nil(t, args.request.Validate()) {
			t.FailNow()
		}
	}
	sharedPostAssert := func(t *testing.T, args args, err error, expectedError string) {
		if expectedError == "" {
			assert.Nil(t, err)
		} else {
			assert.True(t, strings.Contains(err.Error(), expectedError))
		}
	}
	tests := []struct {
		name       string
		args       args
		preAssert  func(t *testing.T, args args)
		postAssert func(t *testing.T, args args, err error, expectedError string)
	}{
		{
			name: "single_container_with_2_points",
			args: args{
				params: &OrderWorkerParams{
					Order: dbo4logist.NewOrderWithData("space1", "order1", mocks4logist.ValidOrderDto1(t)),
				},
				request: dto4logist.AddContainersRequest{
					OrderRequest: dto4logist.NewOrderRequest("space1", "order1"),
					Containers: []dto4logist.NewContainer{
						{
							OrderContainerBase: dbo4logist.OrderContainerBase{
								Type:   "20ft",
								Number: "C111",
							},
							Points: []dto4logist.PointOfNewContainer{
								{
									ShippingPointID: mocks4logist.ShippingPoint1WithSingleContainerID,
									Tasks: []dbo4logist.ShippingPointTask{
										dbo4logist.ShippingPointTaskLoad,
									},
								},
								{
									ShippingPointID: mocks4logist.ShippingPoint1WithSingleContainerID,
									Tasks: []dbo4logist.ShippingPointTask{
										dbo4logist.ShippingPointTaskLoad,
										dbo4logist.ShippingPointTaskUnload,
									},
								},
							},
						},
					},
				},
			},
			preAssert: func(t *testing.T, args args) {
				sharedPreAssert(t, args)
				orderDto := args.params.Order.Dto
				assert.Equal(t, mocks4logist.Dto1ContainerPointsCount, len(orderDto.ContainerPoints))
			},
			postAssert: func(t *testing.T, args args, err error, expectedError string) {
				sharedPostAssert(t, args, err, expectedError)
				orderDto := args.params.Order.Dto
				assert.Equal(t, mocks4logist.Dto1ContainerPointsCount+2, len(orderDto.ContainerPoints))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preAssert(t, tt.args)
			err := addContainersTx(tt.args.params, tt.args.request)
			tt.postAssert(t, tt.args, err, "")
		})
	}
}
