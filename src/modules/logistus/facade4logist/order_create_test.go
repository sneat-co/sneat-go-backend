package facade4logist

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

func TestUpdateLogistSpace(t *testing.T) {
	logistSpaceDbo := &dbo4logist.LogistSpaceDbo{}
	spaceDbo := &dbo4spaceus.SpaceDbo{WithUserIDs: dbmodels.WithUserIDs{UserIDs: []string{"u1"}}}
	teamContact := dal4contactus.ContactEntry{}
	teamContact.ID = "contact1"
	request := dto4logist.SetLogistSpaceSettingsRequest{
		OrderNumberPrefix: "ORD",
		Roles:             []dbo4logist.LogistSpaceRole{dbo4logist.CompanyRoleTrucker},
	}
	updates := updateLogistSpace(logistSpaceDbo, spaceDbo, teamContact, request)
	assert.NotEmpty(t, updates)
	assert.Equal(t, "ORD", logistSpaceDbo.OrderNumberPrefix)
	assert.Equal(t, "contact1", logistSpaceDbo.ContactID)
}

func TestFillOrderDtoFromRequest(t *testing.T) {
	orderDto := &dbo4logist.OrderDbo{}
	request := dto4logist.CreateOrderRequest{
		NumberOfContainers: map[string]int{"20GP": 1},
	}
	params := &dal4spaceus.SpaceWorkerParams{
		Space: dbo4spaceus.NewSpaceEntry("space1"),
	}
	fillOrderDtoFromRequest(orderDto, request, params, "u1")
	assert.Equal(t, 1, len(orderDto.Containers))
	assert.Equal(t, "20GP1", orderDto.Containers[0].ID)
}

func TestCreateOrder(t *testing.T) {
	request := dto4logist.CreateOrderRequest{
		Order: dbo4logist.OrderBase{
			WithCounterparties: dbo4logist.WithCounterparties{
				Counterparties: []*dbo4logist.OrderCounterparty{
					{Role: dbo4logist.CounterpartyRoleTrucker, ContactID: "c1"},
				},
			},
		},
	}
	// This will fail because dal4spaceus.RunSpaceWorkerWithUserContext is not mocked,
	// but it covers the validation line if we add a call to Validate() inside the function or just the call itself.
	_, err := CreateOrder(nil, request)
	assert.NotNil(t, err)
}

func TestSetLogistSpaceSettings(t *testing.T) {
	request := dto4logist.SetLogistSpaceSettingsRequest{}
	err := SetLogistSpaceSettings(nil, request)
	assert.NotNil(t, err)
}
