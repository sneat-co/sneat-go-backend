package facade4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/mocks4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_txSetContainerPointFreightFields(t *testing.T) {
	type args struct {
		params  *OrderWorkerParams
		request dto4logist.SetContainerPointFreightFieldsRequest
	}
	setArgs := func(containerID, shippingPointID, task, name string, value int) args {
		return args{
			request: dto4logist.SetContainerPointFreightFieldsRequest{
				ContainerPointRequest: dto4logist.ContainerPointRequest{
					ContainerID:     containerID,
					ShippingPointID: shippingPointID,
					OrderRequest:    dto4logist.NewOrderRequest("team1", "order1"),
				},
				Task: task,
				Integers: map[string]*int{
					name: &value,
				},
			},
			params: &OrderWorkerParams{
				Order: models4logist.Order{
					Dto: mocks4logist.ValidOrderDto1(t),
				},
			},
		}
	}
	preAssert := func(t *testing.T, args args) {
		containerPoint := args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, args.request.ShippingPointID)
		if !assert.NotNil(t, containerPoint, "container point not found") {
			panic("container point not found")
		}
		if !assert.Equal(t, 1, len(containerPoint.Tasks)) {
			panic("unexpected container point tasks")
		}
		if !assert.Equal(t, models4logist.ShippingPointTaskLoad, containerPoint.Tasks[0]) {
			panic("unexpected container point task")
		}
	}
	postAssertContainerPoint := func(t *testing.T, args args, err error) (*models4logist.ContainerPoint, bool) {
		assert.Nil(t, err)
		assert.True(t, args.params.Changed.ContainerPoints)
		containerPoint := args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, args.request.ShippingPointID)
		if !assert.NotNil(t, containerPoint, "container point not found") {
			return containerPoint, false
		}
		assert.Equal(t, 1, len(containerPoint.Tasks))
		assert.Equal(t, args.request.Task, containerPoint.Tasks[0])
		return containerPoint, true
	}
	tests := []struct {
		name       string
		args       args
		preAssert  func(t *testing.T, args args)
		postAssert func(t *testing.T, args args, err error)
	}{
		{
			name: "set_load_pallets",
			args: setArgs(mocks4logist.Container1ID, mocks4logist.ShippingPoint2With2ContainersID, models4logist.ShippingPointTaskLoad, "numberOfPallets", 3),
			preAssert: func(t *testing.T, args args) {
				preAssert(t, args)
			},
			postAssert: func(t *testing.T, args args, err error) {
				containerPoint, ok := postAssertContainerPoint(t, args, err)
				if !ok {
					return
				}
				assert.Equal(t, *args.request.Integers["numberOfPallets"], containerPoint.ToLoad.NumberOfPallets)
			},
		},
		{
			name: "set_load_weight",
			args: setArgs(mocks4logist.Container1ID, mocks4logist.ShippingPoint2With2ContainersID, models4logist.ShippingPointTaskLoad, "grossWeightKg", 132),
			preAssert: func(t *testing.T, args args) {
				preAssert(t, args)
			},
			postAssert: func(t *testing.T, args args, err error) {
				containerPoint, ok := postAssertContainerPoint(t, args, err)
				if !ok {
					return
				}
				assert.Equal(t, *args.request.Integers["grossWeightKg"], containerPoint.ToLoad.GrossWeightKg)
			},
		},
		{
			name: "set_load_volume",
			args: setArgs(mocks4logist.Container1ID, mocks4logist.ShippingPoint2With2ContainersID, models4logist.ShippingPointTaskLoad, "volumeM3", 4),
			preAssert: func(t *testing.T, args args) {
				preAssert(t, args)
			},
			postAssert: func(t *testing.T, args args, err error) {
				containerPoint, ok := postAssertContainerPoint(t, args, err)
				if !ok {
					return
				}
				assert.Equal(t, *args.request.Integers["volumeM3"], containerPoint.ToLoad.VolumeM3)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preAssert(t, tt.args)
			err := txSetContainerPointFreightFields(tt.args.params, tt.args.request)
			tt.postAssert(t, tt.args, err)
		})
	}
}
