package facade4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/mocks4logist"
	"github.com/stretchr/testify/assert"
)

func Test_deleteSegments(t *testing.T) {
	type args struct {
		params  *OrderWorkerParams
		request dto4logist.DeleteSegmentsRequest
	}
	tests := []struct {
		name         string
		args         args
		preAssert    func(t *testing.T, args args)
		assertResult func(t *testing.T, order *dbo4logist.Order, err error)
	}{
		{
			name: "should_pass",
			args: args{
				params: &OrderWorkerParams{
					Order: dbo4logist.NewOrderWithData("space1", "order1", mocks4logist.ValidOrderDto1(t)),
				},
				request: dto4logist.DeleteSegmentsRequest{
					OrderRequest: dto4logist.NewOrderRequest("space1", "order1"),
					SegmentsFilter: dbo4logist.SegmentsFilter{
						ContainerIDs:      []string{mocks4logist.Container2ID},
						ToShippingPointID: mocks4logist.ShippingPoint1WithSingleContainerID,
					},
				},
			},
			preAssert: func(t *testing.T, args args) {
				assert.Nil(t, args.request.Validate())
				assert.Equal(t, 1, len(args.params.Order.Dto.Segments))
				for _, containerID := range args.request.ContainerIDs {
					var containerPoint *dbo4logist.ContainerPoint
					if args.request.FromShippingPointID != "" {
						containerPoint = args.params.Order.Dto.GetContainerPoint(containerID, args.request.FromShippingPointID)
					} else if args.request.ToShippingPointID != "" {
						containerPoint = args.params.Order.Dto.GetContainerPoint(containerID, args.request.ToShippingPointID)
					}
					assert.NotNilf(t, containerPoint, "ContainerPoints (%d): %v", len(args.params.Order.Dto.ContainerPoints), args.params.Order.Dto.ContainerPoints)
					assert.NotEqual(t, "", containerPoint.Arrival.ScheduledDate)
				}
			},
			assertResult: func(t *testing.T, order *dbo4logist.Order, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 0, len(order.Dto.Segments))
				containerPoint := order.Dto.GetContainerPoint(mocks4logist.Container1ID, mocks4logist.ShippingPoint2With2ContainersID)
				assert.True(t, containerPoint.Arrival == nil || containerPoint.Arrival.ScheduledDate == "", "arrival scheduled date should be empty")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preAssert(t, tt.args)
			err := deleteSegments(tt.args.params, tt.args.request)
			tt.assertResult(t, &tt.args.params.Order, err)
		})
	}
}
