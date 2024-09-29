package models4debtus

import (
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
)

type DebtusUserDbo struct { // TODO: Move back into debtus module
	money.Balanced
	WithTransferCounts
	WithHasDueTransfers
	models4splitus.BillsHolder
}

type DebtusUserEntry = record.DataWithID[string, *DebtusUserDbo]

func NewDebtusUserEntry(userID string) DebtusUserEntry {
	return dal4userus.NewUserModuleEntry(userID, "debtus", new(DebtusUserDbo))
}
