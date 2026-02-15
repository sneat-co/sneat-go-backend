package dbo4logist

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// Order is a context for an order
type Order struct {
	record.WithID[string]
	Dto *OrderDbo
}

// NewOrderKey create new order key
func NewOrderKey(spaceID coretypes.SpaceID, orderID string) *dal.Key {
	if spaceID == "" {
		panic("spaceID is empty")
	}
	if orderID == "" {
		panic("orderID is empty")
	}
	logistSpaceKey := newLogistSpaceKey(spaceID)
	return dal.NewKeyWithParentAndID(logistSpaceKey, OrdersCollection, orderID)
}

// NewOrder creates new order context
func NewOrder(spaceID coretypes.SpaceID, orderID string) (order Order) {
	key := NewOrderKey(spaceID, orderID)
	dto := new(OrderDbo)
	order.ID = orderID
	order.FullID = getOrderFullShortID(spaceID, orderID)
	order.Key = key
	order.Dto = dto
	order.Record = dal.NewRecordWithData(key, dto)
	return
}

func getOrderFullShortID(spaceID coretypes.SpaceID, orderID string) string {
	return string(spaceID) + ":" + orderID
}

func NewOrderWithData(spaceID coretypes.SpaceID, orderID string, dto *OrderDbo) (order Order) {
	key := NewOrderKey(spaceID, orderID)
	order.ID = orderID
	order.Key = key
	order.Dto = dto
	order.Record = dal.NewRecordWithData(key, dto)
	return
}
