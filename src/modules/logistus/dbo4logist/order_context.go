package dbo4logist

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

// Order is a context for an order
type Order struct {
	record.WithID[string]
	Dto *OrderDbo
}

// NewOrderKey create new order key
func NewOrderKey(teamID, orderID string) *dal.Key {
	if teamID == "" {
		panic("spaceID is empty")
	}
	if orderID == "" {
		panic("orderID is empty")
	}
	logistSpaceKey := newLogistSpaceKey(teamID)
	return dal.NewKeyWithParentAndID(logistSpaceKey, OrdersCollection, orderID)
}

// NewOrder creates new order context
func NewOrder(teamID, orderID string) (order Order) {
	key := NewOrderKey(teamID, orderID)
	dto := new(OrderDbo)
	order.ID = orderID
	order.FullID = getOrderFullShortID(teamID, orderID)
	order.Key = key
	order.Dto = dto
	order.Record = dal.NewRecordWithData(key, dto)
	return
}

func getOrderFullShortID(teamID, orderID string) string {
	return teamID + ":" + orderID
}

func NewOrderWithData(teamID, orderID string, dto *OrderDbo) (order Order) {
	key := NewOrderKey(teamID, orderID)
	order.ID = orderID
	order.Key = key
	order.Dto = dto
	order.Record = dal.NewRecordWithData(key, dto)
	return
}
