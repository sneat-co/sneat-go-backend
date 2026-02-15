package facade4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/mocks4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/slice"
)

func Test_txSetContainerPointTask(t *testing.T) {
	type args struct {
		params  *OrderWorkerParams
		request dto4logist.SetContainerPointTaskRequest
	}
	setArgs := func(containerID, shippingPointID, task string, value bool) args {
		return args{
			request: dto4logist.SetContainerPointTaskRequest{
				ContainerPointRequest: dto4logist.ContainerPointRequest{
					ContainerID:     containerID,
					ShippingPointID: shippingPointID,
					OrderRequest:    dto4logist.NewOrderRequest("space1", "order1"),
				},
				Task:  task,
				Value: value,
			},
			params: &OrderWorkerParams{
				Order: dbo4logist.Order{
					Dto: mocks4logist.ValidOrderDto1(t),
				},
			},
		}
	}
	preAssert := func(t *testing.T, args args) bool {
		err := args.params.Order.Dto.Validate()
		return assert.Nilf(t, err, "unexpected error: %v", err)
	}
	postAssert := func(t *testing.T, args args, err error) bool {
		return assert.Nil(t, args.params.Order.Dto.Validate())
	}
	tests := []struct {
		name       string
		args       args
		preAssert  func(t *testing.T, args args)
		postAssert func(t *testing.T, args args, err error)
	}{
		{
			name: "adds_to_non_existing_container_point",
			args: setArgs(mocks4logist.Container1ID, mocks4logist.ShippingPoint1WithSingleContainerID, dbo4logist.ShippingPointTaskLoad, true),
			preAssert: func(t *testing.T, args args) {
				if !preAssert(t, args) {
					t.Fail()
				}
				if !assert.Nil(t, args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, args.request.ShippingPointID), "unexpected container point") {
					t.Fail()
				}
			},
			postAssert: func(t *testing.T, args args, err error) {
				if !postAssert(t, args, err) {
					t.Fail()
				}
				assert.True(t, args.params.Changed.ContainerPoints)
				containerPoint := args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, args.request.ShippingPointID)
				if !assert.NotNil(t, containerPoint, "container point not found") {
					return
				}
				assert.Equal(t, 1, len(containerPoint.Tasks))
				assert.Equal(t, args.request.Task, containerPoint.Tasks[0])
			},
		},
		{
			name: "adds_to_existing_container_point_second_task",
			args: setArgs(mocks4logist.Container1ID, mocks4logist.ShippingPoint2With2ContainersID, dbo4logist.ShippingPointTaskUnload, true),
			preAssert: func(t *testing.T, args args) {
				if !preAssert(t, args) {
					t.Fail()
				}
				containerPoint := args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, args.request.ShippingPointID)
				if !assert.NotNil(t, containerPoint, "container point not found") {
					panic("container point not found")
				}
				if !assert.Equal(t, 1, len(containerPoint.Tasks)) {
					panic("unexpected container point tasks")
				}
				if !assert.Equal(t, dbo4logist.ShippingPointTaskLoad, containerPoint.Tasks[0]) {
					panic("unexpected container point task")
				}
			},
			postAssert: func(t *testing.T, args args, err error) {
				if !postAssert(t, args, err) {
					t.Fail()
				}
				containerPoint := args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, args.request.ShippingPointID)
				if !assert.NotNil(t, containerPoint, "container point not found") {
					return
				}
				assert.Equal(t, 2, len(containerPoint.Tasks))
				assert.True(t, slice.Index(containerPoint.Tasks, dbo4logist.ShippingPointTaskLoad) >= 0)
				assert.True(t, slice.Index(containerPoint.Tasks, dbo4logist.ShippingPointTaskUnload) >= 0)
				assert.Equal(t, dbo4logist.ShippingPointTaskUnload, containerPoint.Tasks[1])
			},
		},
		{
			name: "try_to_adds_to_existing_container_point_existing_task",
			args: setArgs(mocks4logist.Container1ID, mocks4logist.ShippingPoint2With2ContainersID, dbo4logist.ShippingPointTaskLoad, true),
			preAssert: func(t *testing.T, args args) {
				if !preAssert(t, args) {
					t.Fail()
				}
				containerPoint := args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, args.request.ShippingPointID)
				if !assert.NotNil(t, containerPoint, "container point not found") {
					panic("container point not found")
				}
				if !assert.Equal(t, 1, len(containerPoint.Tasks)) {
					panic("unexpected container point tasks")
				}
				if !assert.Equal(t, dbo4logist.ShippingPointTaskLoad, containerPoint.Tasks[0]) {
					panic("unexpected container point task")
				}
			},
			postAssert: func(t *testing.T, args args, err error) {
				if !postAssert(t, args, err) {
					t.Fail()
				}
				assert.False(t, args.params.Changed.HasChanges())
			},
		},
		{
			name: "remove_task_from_existing_container_point",
			args: setArgs(mocks4logist.Container2ID, mocks4logist.ShippingPoint2With2ContainersID, dbo4logist.ShippingPointTaskUnload, false),
			preAssert: func(t *testing.T, args args) {
				if !preAssert(t, args) {
					t.Fail()
				}
				containerPoint := args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, args.request.ShippingPointID)
				if !assert.NotNil(t, containerPoint, "container point not found") {
					panic("container point not found")
				}
				if !assert.Equal(t, 2, len(containerPoint.Tasks)) {
					panic("unexpected container point tasks")
				}
				if !assert.True(t, slice.Index(containerPoint.Tasks, dbo4logist.ShippingPointTaskLoad) >= 0) {
					panic("missing load task")
				}
				if !assert.True(t, slice.Index(containerPoint.Tasks, dbo4logist.ShippingPointTaskUnload) >= 0) {
					panic("missing unload task")
				}
			},
			postAssert: func(t *testing.T, args args, err error) {
				if !postAssert(t, args, err) {
					t.Fail()
				}
				containerPoint := args.params.Order.Dto.GetContainerPoint(args.request.ContainerID, args.request.ShippingPointID)
				if !assert.NotNil(t, containerPoint, "container point not found") {
					return
				}
				assert.Equal(t, 1, len(containerPoint.Tasks))
				assert.Equal(t, dbo4logist.ShippingPointTaskLoad, containerPoint.Tasks[0])
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preAssert(t, tt.args)
			err := txSetContainerPointTask(tt.args.params, tt.args.request)
			tt.postAssert(t, tt.args, err)
		})
	}
}

func TestSetContainerPointTask(t *testing.T) {
	origRunOrderWorker := RunOrderWorker
	defer func() { RunOrderWorker = origRunOrderWorker }()

	RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
		return worker(ctx, nil, &OrderWorkerParams{
			Order: dbo4logist.Order{Dto: &dbo4logist.OrderDbo{}},
		})
	}

	request := dto4logist.SetContainerPointTaskRequest{
		ContainerPointRequest: dto4logist.ContainerPointRequest{
			OrderRequest:    dto4logist.NewOrderRequest("space1", "order1"),
			ContainerID:     "c1",
			ShippingPointID: "sp1",
		},
	}
	err := SetContainerPointTask(nil, request)
	assert.Nil(t, err)
}
