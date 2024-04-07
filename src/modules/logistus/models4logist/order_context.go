package models4logist

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

// Order is a context for an order
type Order struct {
	record.WithID[string]
	Dto *OrderDto
}

// NewOrderKey create new order key
func NewOrderKey(teamID, orderID string) *dal.Key {
	if teamID == "" {
		panic("teamID is empty")
	}
	if orderID == "" {
		panic("orderID is empty")
	}
	logistTeamKey := newLogistTeamKey(teamID)
	return dal.NewKeyWithParentAndID(logistTeamKey, OrdersCollection, orderID)
}

// NewOrder creates new order context
func NewOrder(teamID, orderID string) (order Order) {
	key := NewOrderKey(teamID, orderID)
	dto := new(OrderDto)
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

func NewOrderWithData(teamID, orderID string, dto *OrderDto) (order Order) {
	key := NewOrderKey(teamID, orderID)
	order.ID = orderID
	order.Key = key
	order.Dto = dto
	order.Record = dal.NewRecordWithData(key, dto)
	return
}
