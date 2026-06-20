package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCountryKey(t *testing.T) {
	assert.Equal(t, "country=US", getCountryKey("US"))
}

func TestGetRefNumberKey(t *testing.T) {
	assert.Equal(t, "", getRefNumberKey(""))
	assert.Equal(t, "refNumber=R1", getRefNumberKey("R1"))
}

func TestOrderDbo_UpdateKeys(t *testing.T) {
	o := &OrderDbo{}
	o.Counterparties = []*OrderCounterparty{
		{ContactID: "c1", CountryID: "US", RefNumber: "R1"},
		{ContactID: "c2", CountryID: "GB"},
	}
	o.UpdateKeys()
	assert.Contains(t, o.Keys, getContactKey("c1"))
	assert.Contains(t, o.Keys, getCountryKey("US"))
	assert.Contains(t, o.Keys, getRefNumberKey("R1"))
	assert.Contains(t, o.Keys, getContactKey("c2"))
	assert.Contains(t, o.Keys, getCountryKey("GB"))
	// combined country&contact key
	assert.Contains(t, o.Keys, "country=US&contact=c1")
}
