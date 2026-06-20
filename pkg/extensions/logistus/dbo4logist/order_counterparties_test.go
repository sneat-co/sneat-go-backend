package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func cParty(contactID, role string) *OrderCounterparty {
	return &OrderCounterparty{ContactID: contactID, Role: role, CountryID: "US", Title: "T"}
}

func sampleCounterparties() WithCounterparties {
	return WithCounterparties{Counterparties: []*OrderCounterparty{
		cParty("c1", CounterpartyRoleBuyer),
		cParty("c2", CounterpartyRoleTrucker),
		cParty("c3", CounterpartyRoleTrucker),
	}}
}

func TestWithCounterparties_GetCounterpartyByRole(t *testing.T) {
	v := sampleCounterparties()
	i, c := v.GetCounterpartyByRole(CounterpartyRoleBuyer)
	assert.Equal(t, 0, i)
	assert.Equal(t, "c1", c.ContactID)

	i, c = v.GetCounterpartyByRole(CounterpartyRoleShip)
	assert.Equal(t, -1, i)
	assert.Nil(t, c)
}

func TestWithCounterparties_GetCounterpartiesByRole(t *testing.T) {
	v := sampleCounterparties()
	truckers := v.GetCounterpartiesByRole(CounterpartyRoleTrucker)
	assert.Len(t, truckers, 2)
	assert.Empty(t, v.GetCounterpartiesByRole(CounterpartyRoleShip))
}

func TestWithCounterparties_GetCounterpartiesByContactID(t *testing.T) {
	v := WithCounterparties{Counterparties: []*OrderCounterparty{
		cParty("c1", CounterpartyRoleBuyer),
		cParty("c1", CounterpartyRoleReceiver),
		cParty("c2", CounterpartyRoleTrucker),
	}}
	got := v.GetCounterpartiesByContactID("c1")
	assert.Len(t, got, 2)
	assert.Empty(t, v.GetCounterpartiesByContactID("missing"))
}

func TestWithCounterparties_GetCounterpartyByRoleAndContactID(t *testing.T) {
	v := sampleCounterparties()
	i, c := v.GetCounterpartyByRoleAndContactID(CounterpartyRoleTrucker, "c3")
	assert.Equal(t, 2, i)
	assert.Equal(t, "c3", c.ContactID)

	i, c = v.GetCounterpartyByRoleAndContactID(CounterpartyRoleTrucker, "missing")
	assert.Equal(t, -1, i)
	assert.Nil(t, c)
}

func TestWithCounterparties_GetCounterpartyByContactID(t *testing.T) {
	v := sampleCounterparties()
	i, c := v.GetCounterpartyByContactID("c2")
	assert.Equal(t, 1, i)
	assert.Equal(t, "c2", c.ContactID)

	i, c = v.GetCounterpartyByContactID("missing")
	assert.Equal(t, -1, i)
	assert.Nil(t, c)
}

func TestWithCounterparties_Updates(t *testing.T) {
	u := sampleCounterparties().Updates()
	assert.Len(t, u, 1)
	assert.Equal(t, "counterparties", u[0].FieldName())
}
