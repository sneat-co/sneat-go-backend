package dto4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ValidOrderRequest() OrderRequest {
	return OrderRequest{
		TeamRequest: dto4teamus.ValidTeamRequest(),
		OrderID:     "test-order-id",
	}
}

func TestOrderRequest_Validate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.NotNil(t, OrderRequest{}.Validate())
	})
	t.Run("valid", func(t *testing.T) {
		assert.Nil(t, ValidOrderRequest().Validate())
	})
	t.Run("no_order_id", func(t *testing.T) {
		assert.NotNil(t, OrderRequest{TeamRequest: dto4teamus.ValidTeamRequest()}.Validate())
	})
	t.Run("no_team_id", func(t *testing.T) {
		assert.NotNil(t, OrderRequest{OrderID: "test-order-id"}.Validate())
	})
}
