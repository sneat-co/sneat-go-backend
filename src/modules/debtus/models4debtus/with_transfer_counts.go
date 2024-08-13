package models4debtus

import "github.com/dal-go/dalgo/dal"

type WithTransferCounts struct {

	// TODO: Document purpose and usage
	TransfersWithInterestCount          int `json:"transfersWithInterestCount" firestore:"transfersWithInterestCount,noindex"`
	CountOfAckTransfersByUser           int `json:"countOfAckTransfersByUser,omitempty" firestore:"countOfAckTransfersByUser,noindex,omitempty"` // Do not remove, need for hiding balance/history menu in Telegram
	CountOfReceiptsCreated              int `json:"countOfReceiptsCreated,omitempty" firestore:"countOfReceiptsCreated,noindex,omitempty"`
	CountOfAckTransfersByCounterparties int `json:"countOfAckTransfersByCounterparties,omitempty" firestore:"countOfAckTransfersByCounterparties,noindex,omitempty"` // Do not remove, need for hiding balance/history menu in Telegram
}

func (v *WithTransferCounts) Validate() error {
	return nil
}

type WithHasDueTransfers struct {
	// TODO: Check if we really need this prop and if yes document why
	HasDueTransfers bool `json:"hasDueTransfers,omitempty" firestore:"hasDueTransfers,noindex,omitempty"`
}

func (v *WithHasDueTransfers) SetHasDueTransfers(value bool) (update dal.Update) {
	v.HasDueTransfers = value
	return dal.Update{Field: "hasDueTransfers", Value: value}
}
