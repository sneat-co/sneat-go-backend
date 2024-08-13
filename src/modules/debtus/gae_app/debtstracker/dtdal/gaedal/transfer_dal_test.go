package gaedal

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"testing"
)

func TestNewTransferKey(t *testing.T) {
	const transferID = "12345"
	testStrKey(t, transferID, models4debtus.NewTransferKey(transferID))
}
