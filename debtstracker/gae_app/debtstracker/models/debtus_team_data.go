package models

import "github.com/crediterra/money"

type DebtusTeamData struct {
	money.Balanced
	billsHolder
	TransfersWithInterestCount int  `datastore:",noindex"`
	HasDueTransfers            bool `datastore:",noindex"` // TODO: Check if we really need this prop and if yes document why

	CountOfAckTransfersByUser           int `datastore:",noindex,omitempty"` // Do not remove, need for hiding balance/history menu in Telegram
	CountOfReceiptsCreated              int `datastore:",noindex,omitempty"`
	CountOfAckTransfersByCounterparties int `datastore:",noindex,omitempty"` // Do not remove, need for hiding balance/history menu in Telegram

	//
	BillsCountActive int    `datastore:",noindex,omitempty"`
	BillsJsonActive  string `datastore:",noindex,omitempty"`
	//
	BillSchedulesCountActive int    `datastore:",noindex,omitempty"`
	BillSchedulesJsonActive  string `datastore:",noindex,omitempty"`
}
