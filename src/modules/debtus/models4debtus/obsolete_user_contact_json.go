package models4debtus

//go:generate ffjson $GOFILE

import (
	"fmt"
	"github.com/strongo/validation"
	"time"

	"context"
	"github.com/crediterra/money"
	"github.com/strongo/decimal"
)

type LastTransfer struct {
	ID string    `firestore:"LastTransferID,omitempty" json:"LastTransferID,omitempty"`
	At time.Time `firestore:"LastTransferAt,omitempty" json:"LastTransferAt,omitempty"`
}

type TransferWithInterestJson struct {
	TransferInterest
	TransferID string              `json:"transferID" firestore:"transferID"`
	Direction  TransferDirection   `json:"direction" firestore:"direction"`
	Starts     time.Time           `json:"starts" firestore:"starts"`
	Currency   money.CurrencyCode  `json:"currency,omitempty" firestore:"currency,omitempty"` // TODO: will be obsolete once we group outstanding by currency
	Amount     decimal.Decimal64p2 `json:"amount" firestore:"amount"`
	Returns    TransferReturns     `json:"returns,omitempty" firestore:"returns,omitempty"`
}

func (v TransferWithInterestJson) Validate() error {
	if v.TransferID == "" {
		return validation.NewErrRecordIsMissingRequiredField("transferID")
	}
	if v.Starts.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("starts")
	}
	if v.Currency == "" {
		return validation.NewErrRecordIsMissingRequiredField("currency")
	}
	if v.Amount == 0 {
		return validation.NewErrRecordIsMissingRequiredField("amount")
	}
	if err := v.TransferInterest.Validate(); err != nil {
		return err
	}
	for i, r := range v.Returns {
		if err := r.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("returs[%d]", i), err.Error())
		}
	}
	return nil
}

func (t TransferWithInterestJson) Equal(t2 TransferWithInterestJson) bool {
	if t.TransferInterest != t2.TransferInterest ||
		t.TransferID != t2.TransferID ||
		t.Starts != t2.Starts ||
		t.Currency != t2.Currency ||
		t.Amount != t2.Amount ||
		len(t.Returns) != len(t2.Returns) {
		return false
	}
	for i, r := range t.Returns {
		if r != t2.Returns[i] {
			return false
		}
	}
	return true
}

var _ = TransferInterestCalculable(TransferWithInterestJson{})

func (t TransferWithInterestJson) GetLendingValue() decimal.Decimal64p2 {
	return t.Amount
}

func (t TransferWithInterestJson) GetStartDate() time.Time {
	return t.Starts
}

func (t TransferWithInterestJson) GetReturns() (returns TransferReturns) {
	return t.Returns
}

type UserContactTransfersInfo struct {
	Count                   int                        `json:"count,omitempty" firestore:"count,omitempty"`
	Last                    LastTransfer               `json:"last,omitempty" firestore:"last,omitempty"`
	OutstandingWithInterest []TransferWithInterestJson `json:"outstandingWithInterest,omitempty" firestore:"outstandingWithInterest,omitempty"`
}

func (v *UserContactTransfersInfo) Validate() error {
	if v.Count < 0 {
		return validation.NewErrBadRecordFieldValue("Count", "cannot be negative")
	}
	return nil
}
func (v *UserContactTransfersInfo) Equal(o2 *UserContactTransfersInfo) bool {
	if o2 == nil || v.Count != o2.Count || v.Last != o2.Last || len(v.OutstandingWithInterest) != len(o2.OutstandingWithInterest) {
		return false
	}
	for i, t := range v.OutstandingWithInterest {
		if !t.Equal(o2.OutstandingWithInterest[i]) {
			return false
		}
	}
	return true
}

type DebtusContactStatus string

const (
	DebtusContactStatusActive   DebtusContactStatus = "active"
	DebtusContactStatusArchived DebtusContactStatus = "archived"
)

type DebtusContactBrief struct {
	Status    DebtusContactStatus       `json:"status" firestore:"status"` // We might want to hide a contact from debts list without archiving the contact itself
	Balance   money.Balance             `json:"balance,omitempty"`
	Transfers *UserContactTransfersInfo `json:"transfers,omitempty"`
}

type DebtusContactBriefWithContactID struct {
	ContactID string
	DebtusContactBrief
}

func (o *DebtusContactBrief) Equal(o2 *DebtusContactBrief) bool {
	return o.Status == o2.Status &&
		o.Balance.Equal(o2.Balance) &&
		(o.Transfers == nil && o2.Transfers == nil || o.Transfers != nil && o2.Transfers != nil && o.Transfers.Equal(o2.Transfers))
}

func (o *DebtusContactBrief) BalanceWithInterest(_ context.Context, periodEnds time.Time) (balance money.Balance, err error) {
	balance = make(money.Balance, len(o.Balance))
	for currency, amount := range o.Balance {
		balance[currency] = amount
	}
	if o.Transfers != nil {
		if err = updateBalanceWithInterest(false, balance, o.Transfers.OutstandingWithInterest, periodEnds); err != nil {
			return
		}
	}
	return
}

func NewDebtusContactJson(status DebtusContactStatus, balanced money.Balanced) *DebtusContactBrief {
	result := DebtusContactBrief{
		Status:  status,
		Balance: balanced.Balance,
	}
	if balanced.CountOfTransfers != 0 {
		if balanced.LastTransferID == "" {
			panic(fmt.Sprintf("balanced.CountOfTransfers:%v != 0 && balanced.LastTransferID == 0", balanced.CountOfTransfers))
		}
		if balanced.LastTransferAt.IsZero() {
			panic(fmt.Sprintf("balanced.CountOfTransfers:%v != 0 && balanced.LastTransferAt.IsZero():true", balanced.CountOfTransfers))
		}
		result.Transfers = &UserContactTransfersInfo{
			Count: balanced.CountOfTransfers,
			Last: LastTransfer{
				ID: balanced.LastTransferID,
				At: balanced.LastTransferAt,
			},
		}
	}
	return &result
}
