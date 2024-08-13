package models4debtus

import (
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
)

type DebtusUserDbo struct {
	money.Balanced
	WithTransferCounts
	WithHasDueTransfers
	models4splitus.BillsHolder
}

type DebtusUserEntry = record.DataWithID[string, *DebtusUserDbo]

func NewDebtusUserEntry(userID string) DebtusUserEntry {
	return dal4userus.NewUserModuleEntry(userID, const4debtus.ModuleID, new(DebtusUserDbo))
}
