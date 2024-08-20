package models4debtus

import "github.com/dal-go/dalgo/dal"

type WithTransferCounts struct {

	// TODO: Document purpose and usage
	TransfersWithInterestCount          int `json:"transfersWithInterestCount,omitempty" firestore:"transfersWithInterestCount,omitempty"`
	CountOfAckTransfersByUser           int `json:"countOfAckTransfersByUser,omitempty" firestore:"countOfAckTransfersByUser,omitempty"` // Do not remove, need for hiding balance/history menu in Telegram
	CountOfReceiptsCreated              int `json:"countOfReceiptsCreated,omitempty" firestore:"countOfReceiptsCreated,omitempty"`
	CountOfAckTransfersByCounterparties int `json:"countOfAckTransfersByCounterparties,omitempty" firestore:"countOfAckTransfersByCounterparties,omitempty"` // Do not remove, need for hiding balance/history menu in Telegram
}

func (v *WithTransferCounts) Validate() error {
	return nil
}

type WithHasDueTransfers struct {
	// TODO: Check if we really need this prop and if yes document why
	HasDueTransfers bool `json:"hasDueTransfers,omitempty" firestore:"hasDueTransfers,omitempty"`
}

func (v *WithHasDueTransfers) SetHasDueTransfers(value bool) (update dal.Update) {
	v.HasDueTransfers = value
	return dal.Update{Field: "hasDueTransfers", Value: value}
}
