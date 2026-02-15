package dbo4logist

import (
	"testing"
)

func TestShippingDoc(t *testing.T) {
	if len(ShippingDoc.Fields) == 0 {
		t.Error("expected ShippingDoc to have fields")
	}
	foundNumber := false
	for _, field := range ShippingDoc.Fields {
		if field.ID == "number" {
			foundNumber = true
			break
		}
	}
	if !foundNumber {
		t.Error("expected ShippingDoc to have 'number' field")
	}
}
