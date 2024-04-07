package mocks4logist

import (
	"testing"
)

const (
	orderMustBeValid = "this test order must be valid, got unexpected error: %v"
)

func TestValidEmptyOrder(t *testing.T) {
	order := ValidEmptyOrder(t)
	if err := order.Validate(); err != nil {
		t.Errorf(orderMustBeValid, err)
	}
}

func TestValidOrderWith3UnassignedContainers(t *testing.T) {
	order := ValidOrderWith3UnassignedContainers(t)
	if err := order.Validate(); err != nil {
		t.Errorf(orderMustBeValid, err)
	}
}

func TestValidOrderWith4containers2points1port(t *testing.T) {
	order := ValidOrderDto1(t)
	if err := order.Validate(); err != nil {
		t.Errorf(orderMustBeValid, err)
	}
}
