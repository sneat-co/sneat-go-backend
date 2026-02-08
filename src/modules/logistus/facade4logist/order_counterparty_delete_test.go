package facade4logist

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/stretchr/testify/assert"
)

func TestDeleteOrderCounterparty(t *testing.T) {
	origRunOrderWorker := RunOrderWorker
	defer func() { RunOrderWorker = origRunOrderWorker }()

	RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
		return worker(ctx, nil, &OrderWorkerParams{Order: dbo4logist.Order{Dto: &dbo4logist.OrderDbo{}}})
	}

	request := dto4logist.DeleteOrderCounterpartyRequest{
		OrderRequest: dto4logist.NewOrderRequest("space1", "order1"),
		Role:         dbo4logist.CounterpartyRoleTrucker,
		ContactID:    "c1",
	}
	err := DeleteOrderCounterparty(nil, request)
	assert.Nil(t, err)
}

func Test_deleteOrderCounterpartyTxWorker(t *testing.T) {
	orderDto := &dbo4logist.OrderDbo{
		WithOrderContacts: dbo4logist.WithOrderContacts{
			Contacts: []*dbo4logist.OrderContact{
				{ID: "contact1"},
			},
		},
		OrderBase: dbo4logist.OrderBase{
			WithCounterparties: dbo4logist.WithCounterparties{
				Counterparties: []*dbo4logist.OrderCounterparty{
					{Role: dbo4logist.CounterpartyRoleTrucker, ContactID: "contact1"},
				},
			},
		},
	}
	params := &OrderWorkerParams{
		Order: dbo4logist.Order{Dto: orderDto},
	}

	request := dto4logist.DeleteOrderCounterpartyRequest{
		Role:      dbo4logist.CounterpartyRoleTrucker,
		ContactID: "contact1",
	}

	err := deleteOrderCounterpartyTxWorker(params, request)
	assert.Nil(t, err)
	assert.Empty(t, orderDto.Counterparties)
	assert.Empty(t, orderDto.Contacts)
	assert.True(t, params.Changed.Counterparties)
	assert.True(t, params.Changed.Contacts)
}
