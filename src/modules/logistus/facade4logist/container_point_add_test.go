package facade4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/mocks4logist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_txAddContainerPoints(t *testing.T) {
	type args struct {
		request dto4logist.AddContainerPointsRequest
		params  *OrderWorkerParams
	}

	sharedPreAssert := func(t *testing.T, args args) {
		if !assert.Nil(t, args.request.Validate()) {
			t.Fail()
		}
	}

	sharedPostAssert := func(t *testing.T, args args, err error, wantErr bool) {
		if wantErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}

	tests := []struct {
		name       string
		args       args
		preAssert  func(t *testing.T, args args)
		postAssert func(t *testing.T, args args, err error)
	}{
		{
			name: "adds_single_container_point_with_single_task",
			args: args{
				request: dto4logist.AddContainerPointsRequest{
					OrderRequest: dto4logist.NewOrderRequest("space1", "order1"),
					ContainerPoints: []dbo4logist.ContainerPoint{
						{
							ContainerID:     mocks4logist.Container1ID,
							ShippingPointID: mocks4logist.ShippingPoint1WithSingleContainerID,
							ShippingPointBase: dbo4logist.ShippingPointBase{
								Status: "pending",
								FreightPoint: dbo4logist.FreightPoint{
									Tasks: []string{"load"},
								},
							},
						},
					},
				},
				params: &OrderWorkerParams{
					Order: dbo4logist.NewOrderWithData("space1", "order1", mocks4logist.ValidOrderDto1(t)),
				},
			},
			preAssert: func(t *testing.T, args args) {
				sharedPreAssert(t, args)
				assert.Equal(t, mocks4logist.Dto1ContainerPointsCount, len(args.params.Order.Dto.ContainerPoints))
			},
			postAssert: func(t *testing.T, args args, err error) {
				sharedPostAssert(t, args, err, false)
				assert.Equal(t, mocks4logist.Dto1ContainerPointsCount+1, len(args.params.Order.Dto.ContainerPoints))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preAssert(t, tt.args)
			err := txAddContainerPoints(tt.args.request, tt.args.params)
			tt.postAssert(t, tt.args, err)
		})
	}
}
